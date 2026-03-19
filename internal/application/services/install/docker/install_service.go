package docker

import (
	"synctl/internal/application/interfaces/install"
	install_docker "synctl/internal/application/interfaces/install/docker"
	"synctl/internal/domain"
)

type InstallService struct {
	detector  install.Detector
	inspector install.Inspector
	installer install_docker.DockerInstaller
}

func NewInstallService(
	detector install.Detector,
	inspector install.Inspector,
	installer install_docker.DockerInstaller,
) *InstallService {
	return &InstallService{
		detector:  detector,
		inspector: inspector,
		installer: installer,
	}
}

func (s *InstallService) Install() (*domain.Snapshot, error) {
	status, err := s.detector.Detect()

	if err != nil {
		return nil, err
	}

	if !status.DockerInstalled {
		if err := s.installer.InstallDocker(); err != nil {
			return nil, err
		}
	}

	if !status.DockerRunning {
		if err := s.installer.StartDocker(); err != nil {
			return nil, err
		}
	}

	return s.inspector.Snapshot()

}
