package install_k8s

type K8sInstaller interface {
	InstallCluster() error
	ConfigureCluster() error
	InstallKubectl() error
	InstallK3dCluster() error
	InstallK3dBinary() error
}
