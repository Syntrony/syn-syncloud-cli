package k8s

import (
	"fmt"
	"synctl/internal/application/interfaces"
	dto "synctl/internal/domain/dto"
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

	fmt.Println("Detecting Kubernetes environment...")
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
