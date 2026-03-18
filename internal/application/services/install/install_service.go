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

	state := s.builder.Build(snapshot, nodes)

	s.logger.Info("Syncloud Platform installed successfully!!!")

	return s.repo.Save(state)
}
