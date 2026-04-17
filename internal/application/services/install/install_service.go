package install

import (
	"synctl/internal/application/interfaces"
	run_install "synctl/internal/application/interfaces/install"
	"synctl/internal/application/interfaces/system"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
)

type InstallService struct {
	repo          interfaces.Repository
	runtimes      []run_install.RuntimeInstaller
	nodeInspector system.Inspector
	builder       *state.StateBuilder
	logger        interfaces.Logger
}

func NewInstallService(
	repo interfaces.Repository,
	runtimes []run_install.RuntimeInstaller,
	nodeInspector system.Inspector,
	builder *state.StateBuilder,
	logger interfaces.Logger,
) *InstallService {
	return &InstallService{
		repo:          repo,
		runtimes:      runtimes,
		nodeInspector: nodeInspector,
		builder:       builder,
		logger:        logger,
	}
}

func (s *InstallService) Install() error {
	snapshot := &domain.Snapshot{}

	for _, runtime := range s.runtimes {
		snap, err := runtime.Install()

		if err != nil {
			return err
		}

		snapshot.Resources = append(snapshot.Resources, snap.Resources...)
		snapshot.Records = append(snapshot.Records, snap.Records...)
	}

	nodes, err := s.nodeInspector.Inspect()
	if err != nil {
		return err
	}

	// Preserve existing cluster metadata and node IDs across re-installs
	var existingCluster *domain.Cluster
	if exists, _ := s.repo.Exists(); exists {
		if current, err := s.repo.Load(); err == nil {
			existingCluster = current.Cluster
			existingNodeIndex := make(map[string]domain.Node)
			for _, n := range current.Nodes {
				existingNodeIndex[n.Hostname] = n
			}
			for idx, n := range nodes {
				if prev, ok := existingNodeIndex[n.Hostname]; ok {
					nodes[idx].Id = prev.Id
				}
			}
		}
	}

	state := s.builder.Build(snapshot, existingCluster, nodes)

	return s.repo.Save(state)
}
