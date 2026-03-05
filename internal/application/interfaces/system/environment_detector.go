package system

type EnvironmentDetector interface {
	IsContainer() bool
}
