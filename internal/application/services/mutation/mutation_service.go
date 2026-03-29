package mutation

import (
	"synctl/internal/application/interfaces"
	apply "synctl/internal/application/interfaces/common"
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

func (m *MutationService) ExecuteResource(resources []*domain.Resource, mut mutation.StateMutation) error {
	return m.Execute(resources, mut)
}

func (m *MutationService) ExecuteFile(file string, mut mutation.StateMutation) error {
	resources, err := m.parser.Parse(file)

	if err != nil {
		return err
	}

	return m.Execute(resources, mut)
}

func (m *MutationService) Execute(resources []*domain.Resource, mut mutation.StateMutation) error {
	for _, r := range resources {
		if err := m.validator.Validate(r); err != nil {
			return err
		}
	}

	if err := m.validator.ValidateSet(resources); err != nil {
		return err
	}

	snapshot, err := mut.Mutate(resources)

	if err != nil {
		return err
	}

	var snap domain.Snapshot
	snap.Resources = snapshot

	state := m.builder.Build(&snap, nil)

	return m.repo.Save(state)
}
