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
	status.DockerInstalled = err == nil

	if status.DockerInstalled {
		_, err = d.runner.Run("docker", "info")
		status.DockerRunning = err == nil

		status.DockerPermissionsOk = d.UserHasDockerAccess()

		if !status.DockerPermissionsOk {
			d.logger.Info("Docker installed but user lacks socket permissions (not in docker group)")
		}
	}

	d.logger.Info("docker: installed=" + boolStr(status.DockerInstalled) +
		" running=" + boolStr(status.DockerRunning) +
		" permissions=" + boolStr(status.DockerPermissionsOk))

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
	return strings.Contains(output.Stdout, "docker")
}

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
