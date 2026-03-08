package docker

import (
	"fmt"
	"os"

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

	user := os.Getenv("USER")

	_, err := i.sudo.Run("apt-get", "install", "-y", "docker.io")
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
