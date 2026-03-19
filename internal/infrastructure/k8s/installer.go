package k8s

import (
	"fmt"
	"os"
	"runtime"

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

func (i *Installer) InstallKubectl() error {

	i.logger.Info("Installing kubectl...")

	arch := "amd64"

	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	}

	url := fmt.Sprintf(
		"https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/%s/kubectl",
		arch,
	)

	_, err := i.runner.Run(
		"sh",
		"-c",
		fmt.Sprintf("curl -LO %s", url),
	)

	if err != nil {
		return fmt.Errorf("kubectl download failed: %w", err)
	}

	_, err = i.runner.Run("chmod", "+x", "kubectl")

	if err != nil {
		return fmt.Errorf("kubectl chmod failed: %w", err)
	}

	_, err = i.sudo.Run("mv", "kubectl", "/usr/local/bin/")

	if err != nil {
		return fmt.Errorf("kubectl install failed: %w", err)
	}

	i.logger.Info("kubectl installed successfully")

	return nil
}

func (i *Installer) InstallCluster() error {

	i.logger.Info("Installing cluster...")

	switch i.runtime {

	case domain.Container:
		return i.InstallK3d()

	case domain.VPS:
		return i.InstallK3s()

	default:
		return fmt.Errorf("unknown runtime: %s", i.runtime)
	}
}

func (i *Installer) ConfigureCluster() error {

	switch i.runtime {

	case domain.VPS:
		return i.ConfigureK3s()

	case domain.Container:
		return i.ConfigureK3d()

	default:
		return fmt.Errorf("unknown runtime")
	}
}

func (i *Installer) InstallK3s() error {

	i.logger.Info("Installing k3s cluster...")

	_, err := i.runner.Run(
		"sh",
		"-c",
		"curl -sfL https://get.k3s.io | sh -",
	)

	if err != nil {
		return fmt.Errorf("k3s installation failed: %w", err)
	}

	return nil
}

func (i *Installer) InstallK3d() error {

	i.logger.Info("Installing k3d cluster...")

	_, err := i.runner.Run(
		"sh",
		"-c",
		"curl -sfL https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash",
	)

	if err != nil {
		return fmt.Errorf("k3d install failed: %w", err)
	}

	_, err = i.runner.Run("k3d", "cluster", "get", "syncloud")

	if err == nil {
		i.logger.Info("k3d cluster already exists")
		return nil
	}

	_, err = i.runner.Run(
		"k3d",
		"cluster",
		"create",
		"syncloud",
	)

	if err != nil {
		return fmt.Errorf("k3d cluster creation failed: %w", err)
	}

	i.logger.Info("k3d cluster created")

	return nil
}

func (i *Installer) ConfigureK3s() error {

	i.logger.Info("Configuring kubectl for k3s...")

	home := os.Getenv("HOME")

	if home == "" {
		return fmt.Errorf("HOME environment variable not found")
	}

	kubeDir := fmt.Sprintf("%s/.kube", home)
	kubeConfig := fmt.Sprintf("%s/config", kubeDir)

	_, err := i.runner.Run("mkdir", "-p", kubeDir)

	if err != nil {
		return fmt.Errorf("creating kube directory failed: %w", err)
	}

	_, err = i.sudo.Run(
		"cp",
		"/etc/rancher/k3s/k3s.yaml",
		kubeConfig,
	)

	if err != nil {
		return fmt.Errorf("copy kubeconfig failed: %w", err)
	}

	_, err = i.sudo.Run(
		"chown",
		fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		kubeConfig,
	)

	if err != nil {
		return fmt.Errorf("chown kubeconfig failed: %w", err)
	}

	i.logger.Info("kubectl configured")

	return nil
}

func (i *Installer) ConfigureK3d() error {

	i.logger.Info("Configuring kubectl for k3d...")

	_, err := i.runner.Run(
		"k3d",
		"kubeconfig",
		"merge",
		"syncloud",
		"--kubeconfig-switch-context",
	)

	if err != nil {
		return fmt.Errorf("k3d kubeconfig merge failed: %w", err)
	}

	i.logger.Info("kubectl configured")

	return nil
}
