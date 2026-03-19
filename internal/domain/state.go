package domain

type State struct {
	Version   string     `json:"version"`
	Cluster   *Cluster   `json:"cluster,omitempty"`
	Nodes     []Node     `json:"nodes"`
	Resources []Resource `json:"resources"`
	Dns       *Dns       `json:"dns,omitempty"`
}
