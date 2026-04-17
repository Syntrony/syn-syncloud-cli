package docker

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/application/interfaces/install"
	install_docker "synctl/internal/application/interfaces/install/docker"
	"synctl/internal/application/services/common/banner"
	"synctl/internal/domain"
)

type InstallService struct {
	detector  install.Detector
	inspector install.Inspector
	installer install_docker.DockerInstaller
	runner    interfaces.CommandRunner
}

func NewInstallService(
	detector install.Detector,
	inspector install.Inspector,
	installer install_docker.DockerInstaller,
	runner interfaces.CommandRunner,
	_ interfaces.Logger,
) *InstallService {
	return &InstallService{
		detector:  detector,
		inspector: inspector,
		installer: installer,
		runner:    runner,
	}
}

func (s *InstallService) Install() (*domain.Snapshot, error) {
	status, err := s.detector.Detect()
	if err != nil {
		return nil, err
	}

	freshInstall := false
	dockerGroupConfigured := false

	if !status.DockerInstalled {
		if err := s.installer.InstallDocker(); err != nil {
			return nil, err
		}
		freshInstall = true
		dockerGroupConfigured = true
	} else {
		if !status.DockerRunning {
			if err := s.installer.StartDocker(); err != nil {
				return nil, err
			}
		}

		if !status.DockerPermissionsOk {
			user, err := s.runner.Whoami()
			if err != nil {
				return nil, err
			}
			if err := s.installer.ConfigureDockerPermissions(user); err != nil {
				return nil, err
			}
			dockerGroupConfigured = true
		}
	}

	if dockerGroupConfigured {
		banner.PrintSessionRestartWarning()
	}

	// On fresh install docker group membership is not active in the current session.
	// There are no containers on a fresh install anyway, so return empty snapshot.
	if freshInstall {
		return &domain.Snapshot{}, nil
	}

	return s.inspector.Snapshot()
}
