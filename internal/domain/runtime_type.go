package domain

type RuntimeType string

const (
	// K3s       RuntimeType = "k3s"
	// K3d       RuntimeType = "k3d"
	Container RuntimeType = "container"
	VPS       RuntimeType = "vps"
)
