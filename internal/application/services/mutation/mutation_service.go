package mutation

import (
	"synctl/internal/application/interfaces"
	apply "synctl/internal/application/interfaces/common"
	"synctl/internal/application/interfaces/mutation"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
	system "synctl/internal/infrastructure/system"
)

type MutationService struct {
	parser    apply.ResourceParser
	validator apply.Validator
	repo      interfaces.Repository
	builder   *state.StateBuilder
	logger    interfaces.Logger
	inspector *system.Inspector
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

func (m *MutationService) WithInspector(inspector *system.Inspector) *MutationService {
	m.inspector = inspector
	return m
}

func (m *MutationService) getExistingCluster() *domain.Cluster {
	exists, err := m.repo.Exists()
	if err != nil || !exists {
		return nil
	}
	st, err := m.repo.Load()
	if err != nil || st.Cluster == nil {
		return nil
	}
	return st.Cluster
}

func (m *MutationService) getExistingNodes() []domain.Node {
	exists, err := m.repo.Exists()
	if err != nil || !exists {
		return nil
	}
	st, err := m.repo.Load()
	if err != nil || st.Nodes == nil {
		return nil
	}
	return st.Nodes
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

	cluster := m.getExistingCluster()
	nodes := m.getExistingNodes()

	if (nodes == nil || len(nodes) == 0) && m.inspector != nil {
		detectedNodes, err := m.inspector.Inspect()
		if err == nil && len(detectedNodes) > 0 {
			nodes = detectedNodes
		}
	}

	state := m.builder.Build(&snap, cluster, nodes)

	return m.repo.Save(state)
}
