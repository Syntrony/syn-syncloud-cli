package mutation

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/application/interfaces/apply"
	"synctl/internal/application/interfaces/mutation"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
)

type MutationService struct {
	parser    apply.ResourceParser
	validator apply.Validator
	repo      interfaces.Repository
	builder   *state.StateBuilder
	logger    interfaces.Logger
}

func NewMutationService(parser apply.ResourceParser, validator apply.Validator, repo interfaces.Repository, builder *state.StateBuilder, logger interfaces.Logger) *MutationService {
	return &MutationService{
		parser:    parser,
		validator: validator,
		repo:      repo,
		builder:   builder,
		logger:    logger,
	}
}

func (m *MutationService) Execute(file string, mutation mutation.StateMutation) error {
	resources, err := m.parser.Parse(file)

	if err != nil {
		return err
	}

	for _, r := range resources {
		if err := m.validator.Validate(r); err != nil {
			return err
		}
	}

	if err := m.validator.ValidateSet(resources); err != nil {
		return err
	}

	snapshot, err := mutation.Mutate(resources)

	if err != nil {
		return err
	}

	var snap domain.Snapshot
	snap.Resources = snapshot

	state := m.builder.Build(&snap, nil)

	return m.repo.Save(state)
}
