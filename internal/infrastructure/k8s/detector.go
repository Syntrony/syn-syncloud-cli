package k8s

import (
	"synctl/internal/application/interfaces"
	dto "synctl/internal/domain/dto"
)

type Detector struct {
	runner interfaces.CommandRunner
	logger interfaces.Logger
}

func NewDetector(runner interfaces.CommandRunner, logger interfaces.Logger) *Detector {
	return &Detector{
		runner: runner,
		logger: logger,
	}
}

func (d *Detector) Detect() (*dto.Status, error) {

	d.logger.Info("Detecting Kubernetes environment...")
	status := &dto.Status{}

	_, err := d.runner.Run("kubectl", "version", "--client", "-o", "json")

	if err != nil {
		status.KubectlInstalled = false
		return status, nil
	}

	status.KubectlInstalled = true

	_, err = d.runner.Run("kubectl", "cluster-info")

	if err != nil {
		status.ClusterReachable = false
		return status, nil
	}

	status.ClusterReachable = true

	return status, nil
}
