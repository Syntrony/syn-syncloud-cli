package docker

import (
	"fmt"
	"time"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
)

type Installer struct {
	runner  interfaces.CommandRunner
	sudo    interfaces.PrivilegedRunner
	runtime domain.RuntimeType
	logger  interfaces.Logger
}

func NewInstaller(
	runner interfaces.CommandRunner,
	sudo interfaces.PrivilegedRunner,
	runtime domain.RuntimeType,
	logger interfaces.Logger,
) *Installer {
	return &Installer{
		runner:  runner,
		sudo:    sudo,
		runtime: runtime,
		logger:  logger,
	}
}

func (i *Installer) InstallDocker() error {

	if i.runtime == domain.Container {
		i.logger.Info("Skipping docker install inside container runtime")
		return nil
	}

	i.logger.Info("Installing docker...")

	user, err := i.sudo.Whoami()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Use official Docker install script — supports all major Linux distros
	// and avoids the 'docker.io' package availability issue on some systems.
	_, err = i.sudo.Run("sh", "-c", "curl -fsSL https://get.docker.com | sh")
	if err != nil {
		return fmt.Errorf("docker install failed: %w", err)
	}

	_, err = i.sudo.Run("systemctl", "enable", "docker")
	if err != nil {
		return fmt.Errorf("docker enable failed: %w", err)
	}

	_, err = i.sudo.Run("systemctl", "start", "docker")
	if err != nil {
		return fmt.Errorf("docker start failed: %w", err)
	}

	if err := i.waitForDocker(); err != nil {
		i.logger.Info("Docker daemon not ready after install, continuing: " + err.Error())
	}

	return i.ConfigureDocker(user)
}

func (i *Installer) StartDocker() error {

	if i.runtime == domain.Container {
		i.logger.Info("Skipping docker start inside container runtime")
		return nil
	}

	_, err := i.sudo.Run("systemctl", "start", "docker")

	if err != nil {
		return fmt.Errorf("docker initializing failed: %w", err)
	}

	if err := i.waitForDocker(); err != nil {
		i.logger.Info("Docker daemon not ready, continuing: " + err.Error())
	}

	return nil
}

// waitForDocker polls sudo docker info until the daemon is ready or times out.
func (i *Installer) waitForDocker() error {
	i.logger.Info("Waiting for docker daemon to be ready...")
	for attempt := 0; attempt < 15; attempt++ {
		if _, err := i.sudo.Run("docker", "info"); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("docker daemon did not become ready in time")
}

func (i *Installer) ConfigureDockerPermissions(user string) error {
	if err := i.ConfigureDocker(user); err != nil {
		return err
	}
	i.logger.Info("User '" + user + "' added to docker group.")
	return nil
}

func (i *Installer) ConfigureDocker(user string) error {

	if i.runtime == domain.Container {
		_, err := i.runner.Run("usermod", "-aG", "docker", user)

		if err != nil {
			return fmt.Errorf("docker user config failed: %w", err)
		}
		return nil
	}

	_, err := i.sudo.Run("usermod", "-aG", "docker", user)

	if err != nil {
		return fmt.Errorf("docker user config failed: %w", err)
	}

	return nil
}
