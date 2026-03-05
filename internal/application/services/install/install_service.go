package install

import (
	"encoding/json"
	"os"
	run_install "synctl/internal/application/interfaces/install"
	"synctl/internal/application/interfaces/system"
	"synctl/internal/domain"
)

type InstallService struct {
	runtimes      []run_install.RuntimeInstaller
	nodeInspector system.Inspector
	builder       *StateBuilder
}

func NewInstallService(runtimes []run_install.RuntimeInstaller, nodeInspector system.Inspector, builder *StateBuilder) *InstallService {
	return &InstallService{
		runtimes:      runtimes,
		nodeInspector: nodeInspector,
		builder:       builder,
	}
}

func (s *InstallService) Install(mode domain.InstallMode) error {
	snapshot := &domain.Snapshot{}

	for _, runtime := range s.runtimes {
		snap, err := runtime.Install(mode)

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

	jsonData, err := json.MarshalIndent(state, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile("helpers/state.json", jsonData, 0644)
}
