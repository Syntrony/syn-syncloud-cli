package dto

type Status struct {
	KubectlInstalled     bool `json:"kubectlInstalled,omitempty"`
	K3sInstalled         bool `json:"k3sInstalled,omitempty"`
	KubeconfigConfigured bool `json:"kubeconfigConfigured,omitempty"`
	ClusterReachable     bool `json:"clusterReachable,omitempty"`
	Version              string `json:"kustomizeVersion,omitempty"`

	DockerInstalled     bool `json:"dockerInstalled,omitempty"`
	DockerRunning       bool `json:"dockerRunning,omitempty"`
	DockerPermissionsOk bool `json:"dockerPermissionsOk,omitempty"`

	DnsInstalled  bool `json:"dnsInstalled,omitempty"`
	DnsConfigured bool `json:"dnsConfigured,omitempty"`
}
