package executor

import (
	"bytes"
	"os/exec"
	"synctl/internal/domain/dto"
	"synctl/internal/domain/interfaces"
)

type ExecRunner struct {
	Logger interfaces.Logger
}

func NewExecRunner(logger interfaces.Logger) *ExecRunner {
	return &ExecRunner{
		Logger: logger,
	}
}

func (r *ExecRunner) Run(name string, args ...string) (*dto.CommandResult, error) {
	r.Logger.Info("Running command: " + name)

	cmd := exec.Command(name, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &dto.CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		result.Error = err.Error()
		r.Logger.Error("Command execution failed: " + err.Error())
		return result, err
	}

	return result, nil
}
