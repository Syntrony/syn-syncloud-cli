package docker

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"
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
	d.logger.Info("Detecting docker environment...")
	status := &dto.Status{}

	_, err := d.runner.Run("docker", "--version")

	if err != nil {
		status.DockerInstalled = false
		return status, nil
	}

	status.DockerInstalled = true

	_, err = d.runner.Run("docker", "info")

	if err != nil {
		status.DockerRunning = false
		return status, nil
	}

	status.DockerRunning = true

	return status, nil
}
