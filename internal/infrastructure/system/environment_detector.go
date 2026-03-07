package system

import (
	"os"
	"strings"
)

type EnvironmentDetector struct{}

func NewEnvironmentDetector() *EnvironmentDetector {
	return &EnvironmentDetector{}
}

func (e *EnvironmentDetector) IsContainer() bool {
	data, err := os.ReadFile("/proc/1/cgroup")

	if err == nil {
		content := string(data)

		if strings.Contains(content, "docker") ||
			strings.Contains(content, "containerd") ||
			strings.Contains(content, "kubepods") {
			return true
		}
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}

func (e *EnvironmentDetector) IsRoot() bool {
	return os.Geteuid() == 0
}
