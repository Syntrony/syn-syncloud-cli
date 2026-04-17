package install_k8s

type K8sInstaller interface {
	InstallKubectl() error
	InstallCluster() error
	StartCluster() error
	ConfigureCluster() error
	ConfigureKubeconfig() error
}
