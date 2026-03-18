package apply

import (
	"fmt"

	"synctl/internal/application/interfaces"
	iApply "synctl/internal/application/interfaces/apply"
	"synctl/internal/application/services/apply/diff"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
)

type ApplyService struct {
	repo      interfaces.Repository
	builder   *state.StateBuilder
	parser    iApply.ResourceParser
	validator iApply.Validator
	runtimes  []iApply.RuntimeReconciler
	logger    interfaces.Logger
}

func NewApplyService(repo interfaces.Repository, builder *state.StateBuilder, parser iApply.ResourceParser,
	validator iApply.Validator, runtimes []iApply.RuntimeReconciler, logger interfaces.Logger) *ApplyService {
	return &ApplyService{
		repo:      repo,
		builder:   builder,
		parser:    parser,
		validator: validator,
		runtimes:  runtimes,
		logger:    logger,
	}
}

func (s *ApplyService) Apply(file string) error {
	s.logger.Info("Applying Resource...")
	resources, err := s.parser.Parse(file)

	if err != nil {
		return err
	}

	for _, resource := range resources {
		if err := s.validator.Validate(resource); err != nil {
			return err
		}
	}

	if err := s.validator.ValidateSet(resources); err != nil {
		return err
	}

	state, err := s.repo.Load()
	if err != nil {
		return err
	}

	snapshot, err := s.repo.Upsert(resources)

	if err != nil {
		return err
	}

	state = s.builder.Build(snapshot, nil)

	if err := s.repo.Save(state); err != nil {
		return err
	}

	ctx := domain.ReconcileContext{
		Resources: resources,
		State:     state,
	}

	for _, runtime := range s.runtimes {
		var runtimeResources []*domain.Resource

		for _, resource := range ctx.Resources {
			if resource.Runtime == runtime.Runtime() {
				runtimeResources = append(runtimeResources, resource)
			}
		}

		if len(runtimeResources) == 0 {
			continue
		}

		var planned []*domain.Action

		for _, res := range runtimeResources {
			actual, err := runtime.Observe(res)
			if err != nil {
				return fmt.Errorf("observing resource %s: %w", res.Name, err)
			}

			action := diff.Compare(res, actual)
			planned = append(planned, action)
		}

		subCtx := domain.ReconcileContext{
			Resources: runtimeResources,
			State:     ctx.State,
			Actions:   planned,
		}

		if err := runtime.Reconcile(subCtx); err != nil {
			return err
		}
	}

	return nil
}
