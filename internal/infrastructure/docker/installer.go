package docker

import (
	"fmt"

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

	_, err = i.sudo.Run("apt-get", "install", "-y", "docker.io")
	if err != nil {
		return fmt.Errorf("docker install failed: %w", err)
	}

	_, err = i.sudo.Run("systemctl", "enable", "docker")
	if err != nil {
		return fmt.Errorf("docker enable failed: %w", err)
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

	return nil
}

func (i *Installer) ConfigureDockerPermissions(user string) error {
	if err := i.ConfigureDocker(user); err != nil {
		return err
	}
	i.logger.Info("User '" + user + "' added to docker group.")
	i.logger.Info("IMPORTANT: The current session does not reflect this change.")
	i.logger.Info("Run 'newgrp docker' or start a new SSH session to use Docker without sudo.")
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
