package mutation

import (
	"fmt"
	"synctl/internal/application/interfaces"
	apply "synctl/internal/application/interfaces/common"
	"synctl/internal/application/interfaces/mutation"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
	resource "synctl/internal/domain/resource"
	system "synctl/internal/infrastructure/system"
)

type MutationService struct {
	parser    apply.ResourceParser
	validator apply.Validator
	repo      interfaces.Repository
	builder   *state.StateBuilder
	logger    interfaces.Logger
	inspector *system.Inspector
	dryRun    bool
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

func (m *MutationService) AsDryRun() *MutationService {
	m.dryRun = true
	return m
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

func (m *MutationService) getExistingDns() []resource.Record {
	exists, err := m.repo.Exists()
	if err != nil || !exists {
		return nil
	}
	st, err := m.repo.Load()
	if err != nil || st.Dns == nil {
		return nil
	}
	return st.Dns.Records
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

	// Separate DNS resources from regular runtime resources
	var regularResources []*domain.Resource
	var incomingRecords []resource.Record

	for _, r := range resources {
		if r.Runtime == "dns" {
			incomingRecords = append(incomingRecords, resource.Record{
				Name:   r.Name,
				Server: specString(r.Spec, "server"),
			})
		} else {
			regularResources = append(regularResources, r)
		}
	}

	snapshot, err := mut.Mutate(regularResources)

	if err != nil {
		return err
	}

	var snap domain.Snapshot
	snap.Resources = snapshot
	snap.Records = mut.MutateDns(incomingRecords, m.getExistingDns())

	cluster := m.getExistingCluster()
	nodes := m.getExistingNodes()

	if (nodes == nil || len(nodes) == 0) && m.inspector != nil {
		detectedNodes, err := m.inspector.Inspect()
		if err == nil && len(detectedNodes) > 0 {
			nodes = detectedNodes
		}
	}

	state := m.builder.Build(&snap, cluster, nodes)

	if m.dryRun {
		m.logger.Info(fmt.Sprintf("[dry-run] would write %d resource(s) to state — no changes applied", len(state.Resources)))
		return nil
	}

	return m.repo.Save(state)
}

func specString(spec map[string]interface{}, key string) string {
	if spec == nil {
		return ""
	}
	if val, ok := spec[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
