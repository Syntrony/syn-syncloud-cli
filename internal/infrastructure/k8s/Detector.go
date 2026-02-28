package k8s

import (
	dto "synctl/internal/domain/dto"
	"synctl/internal/domain/interfaces"
)

type Detector struct {
	runner interfaces.CommandRunner
}

func NewDetector(runner interfaces.CommandRunner) *Detector {
	return &Detector{
		runner: runner,
	}
}

func (d *Detector) Detect() (*dto.Status, error) {
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
