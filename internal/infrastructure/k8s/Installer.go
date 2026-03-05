package k8s

import (
	"fmt"
	"synctl/internal/application/interfaces"
)

type Installer struct {
	runner interfaces.CommandRunner
	sudo   interfaces.PrivilegedRunner
}

func NewInstaller(runner interfaces.CommandRunner, sudo interfaces.PrivilegedRunner) *Installer {
	return &Installer{
		runner: runner,
		sudo:   sudo,
	}
}

func (i *Installer) InstallKubectl() error {
	fmt.Println("Installing kubectl...")
	_, err := i.runner.Run(
		"sh",
		"-c",
		"curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl",
	)

	if err != nil {
		return err
	}

	_, err = i.runner.Run("chmod", "+x", "kubectl")

	if err != nil {
		return err
	}

	_, err = i.sudo.Run("mv", "kubectl", "/usr/local/bin/")

	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) InstallCluster() error {
	fmt.Println("Installing k3s cluster...")
	_, err := i.runner.Run("sh", "-c", "curl -sfL https://get.k3s.io | sh -")

	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) ConfigureCluster() error {
	fmt.Println("Configuring kubectl to access k3s cluster...")
	// For k3s, kubeconfig is typically located at /etc/rancher/k3s/k3s.yaml
	_, err := i.runner.Run("sh", "-c", "mkdir -p $HOME/.kube")

	if err != nil {
		return err
	}

	_, err = i.sudo.Run("sh", "-c",
		"cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config")

	if err != nil {
		return err
	}

	_, err = i.sudo.Run("sh", "-c",
		"chown $(id -u):$(id -g) $HOME/.kube/config")

	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) InstallK3dCluster() error {
	_, err := i.runner.Run("k3d", "version")

	if err != nil {
		if err := i.InstallK3dBinary(); err != nil {
			return err
		}
	}

	_, err = i.runner.Run(
		"k3d", "cluster", "create", "syncloud-dev", "--agents", "1",
	)

	return err
}

func (i *Installer) InstallK3dBinary() error {

	_, err := i.runner.Run(
		"bash",
		"-c",
		"curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash",
	)

	return err
}
