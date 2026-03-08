package dto

type Status struct {
	KubectlInstalled bool   `json:"kubectlInstalled,omitempty"`
	ClusterReachable bool   `json:"clusterReachable,omitempty"`
	Version          string `json:"kustomizeVersion,omitempty"`

	DockerInstalled bool `json:"dockerInstalled,omitempty"`
	DockerRunning   bool `json:"dockerRunning,omitempty"`
}
