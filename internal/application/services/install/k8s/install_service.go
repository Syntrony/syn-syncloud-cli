package k8s

import (
	install "synctl/internal/application/interfaces/install"
	install_k8s "synctl/internal/application/interfaces/install/k8s"
	"synctl/internal/domain"
)

type InstallService struct {
	detector  install.Detector
	inspector install.Inspector
	installer install_k8s.K8sInstaller
}

func NewInstallService(
	detector install.Detector,
	inspector install.Inspector,
	installer install_k8s.K8sInstaller,
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

	if !status.KubectlInstalled {
		if err := s.installer.InstallKubectl(); err != nil {
			return nil, err
		}
	}

	if !status.ClusterReachable {
		if err := s.installer.InstallCluster(); err != nil {
			return nil, err
		}
		if err := s.installer.ConfigureCluster(); err != nil {
			return nil, err
		}
	}

	return s.inspector.Snapshot()
}
