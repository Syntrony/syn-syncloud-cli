package install_k8s

type K8sInstaller interface {
	InstallKubectl() error
	InstallCluster() error
	ConfigureCluster() error
	ConfigureKubeconfig() error
}
