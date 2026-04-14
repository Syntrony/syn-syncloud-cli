package dto

type StatusDto struct {
	StateFound    bool
	Version       string
	ClusterName   string
	ClusterMode   string
	NodeCount     int
	ResourceCount int
	Docker        string
	Kubernetes    string
}
