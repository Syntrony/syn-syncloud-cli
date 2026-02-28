package k8s

import (
	"synctl/internal/domain"
	"synctl/internal/domain/interfaces"
)

type InstallService struct {
	detector  interfaces.Detector
	inspector interfaces.Inspector
	stateRepo interfaces.Repository
}

// func NewInstallService(
// 	detector interfaces.Detector,
// 	inspector interfaces.Inspector,
// 	stateRepo interfaces.Repository,
// ) *InstallService

func (s *InstallService) Execute() (*domain.Snapshot, error) {
	status, err := s.detector.Detect()
	if err != nil {
		return nil, err
	}

	if !status.KubectlInstalled {
		if err := s.InstallKubectl(); err != nil {
			return nil, err
		}
	}

	if !status.ClusterReachable {
		if err := s.InstallCluster(); err != nil {
			return nil, err
		}
	}

	return s.inspector.Snapshot()
}
