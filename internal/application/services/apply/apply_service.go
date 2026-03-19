package apply

import (
	"synctl/internal/application/interfaces"
	iApply "synctl/internal/application/interfaces/apply"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
)

type ApplyService struct {
	repo      interfaces.Repository
	builder   *state.StateBuilder
	parser    iApply.ResourceParser
	validator iApply.Validator
	logger    interfaces.Logger
}

func NewApplyService(repo interfaces.Repository, builder *state.StateBuilder, parser iApply.ResourceParser,
	validator iApply.Validator, logger interfaces.Logger) *ApplyService {
	return &ApplyService{
		repo:      repo,
		builder:   builder,
		parser:    parser,
		validator: validator,
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

	snapshot, err := s.repo.Upsert(resources)

	if err != nil {
		return err
	}

	var snap domain.Snapshot
	snap.Resources = snapshot

	state := s.builder.Build(&snap, nil)

	s.logger.Info("Resource created")
	return s.repo.Save(state)
}
