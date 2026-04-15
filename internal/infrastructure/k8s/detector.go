package k8s

import (
	"os"
	"path/filepath"

	"synctl/internal/application/interfaces"
	dto "synctl/internal/domain/dto"
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
	d.logger.Info("Detecting Kubernetes environment...")
	status := &dto.Status{}

	// 1. kubectl binary
	_, err := d.runner.Run("kubectl", "version", "--client", "-o", "json")
	status.KubectlInstalled = err == nil

	// 2. K3s instalado — independiente de KUBECONFIG
	_, err = d.runner.Run("k3s", "--version")
	status.K3sInstalled = err == nil

	// 3. KUBECONFIG configurado
	status.KubeconfigConfigured = d.isKubeconfigAvailable()

	// 4. Cluster accesible — solo si KUBECONFIG existe
	if status.KubeconfigConfigured {
		_, err = d.runner.Run("kubectl", "cluster-info")
		status.ClusterReachable = err == nil
	}

	d.logger.Info("kubectl: " + boolStr(status.KubectlInstalled) +
		" | k3s: " + boolStr(status.K3sInstalled) +
		" | kubeconfig: " + boolStr(status.KubeconfigConfigured) +
		" | cluster: " + boolStr(status.ClusterReachable))

	return status, nil
}

func (d *Detector) isKubeconfigAvailable() bool {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		if _, err := os.Stat(kubeconfig); err == nil {
			return true
		}
	}
	home := os.Getenv("HOME")
	if home == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(home, ".kube", "config"))
	return err == nil
}

func boolStr(v bool) string {
	if v {
		return "✓"
	}
	return "✗"
}
