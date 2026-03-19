package docker

import (
	"strings"

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

	if !d.UserHasDockerAccess() {
		status.DockerRunning = false
		d.logger.Info("Docker installed but user lacks permissions (not in docker group)")
	}

	return status, nil
}

func (d *Detector) UserHasDockerAccess() bool {
	user, err := d.runner.Whoami()
	if err != nil || user == "" {
		return false
	}

	output, err := d.runner.Run("groups", user)
	if err != nil {
		return false
	}

	return strings.Contains(string(output.Stdout), "docker")
}
