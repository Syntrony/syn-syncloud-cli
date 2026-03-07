package executor

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"
)

type SudoRunner struct {
	runner interfaces.CommandRunner
}

func NewSudoRunner(runner interfaces.CommandRunner) *SudoRunner {
	return &SudoRunner{runner: runner}
}

func (r *SudoRunner) Run(cmd string, args ...string) (*dto.CommandResult, error) {

	fullArgs := append([]string{cmd}, args...)

	return r.runner.Run("sudo", fullArgs...)
}
