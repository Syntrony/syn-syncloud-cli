package install

import (
	"fmt"
	"synctl/internal/application/interfaces"
	run_install "synctl/internal/application/interfaces/install"
	"synctl/internal/application/interfaces/system"
	"synctl/internal/domain"
)

type InstallService struct {
	repo          interfaces.Repository
	runtimes      []run_install.RuntimeInstaller
	nodeInspector system.Inspector
	builder       *StateBuilder
}

func NewInstallService(repo interfaces.Repository, runtimes []run_install.RuntimeInstaller, nodeInspector system.Inspector, builder *StateBuilder) *InstallService {
	return &InstallService{
		repo:          repo,
		runtimes:      runtimes,
		nodeInspector: nodeInspector,
		builder:       builder,
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
	}

	nodes, err := s.nodeInspector.Inspect()

	if err != nil {
		return err
	}

	state := s.builder.Build(snapshot, nodes)

	fmt.Println("Syncloud Platform installed successfully!!!")

	return s.repo.Save(state)
}
