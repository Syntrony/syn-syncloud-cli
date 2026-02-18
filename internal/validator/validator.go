package validator

import (
	"github.com/Syntrony/syn-sycloud-cli/internal/system"
)

type Result struct {
	DockerInstalled  bool
	KubectlInstalled bool
}

func RunValidator() Result {
	return Result{
		DockerInstalled:  system.CommandExists("docker"),
		KubectlInstalled: system.CommandExists("kubectl"),
	}
}
