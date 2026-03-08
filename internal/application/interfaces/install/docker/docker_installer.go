package docker

type DockerInstaller interface {
	InstallDocker() error
	StartDocker() error
}
