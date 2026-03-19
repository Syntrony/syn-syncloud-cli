package executor

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"
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
	fullCmd := formatCommand(name, args)
	r.Logger.Info("Running command: " + fullCmd)

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
		r.Logger.Error(fmt.Sprintf(
			"Command failed: %s\nstderr: %s",
			fullCmd,
			result.Stderr,
		))
		return result, err
	}
	// r.Logger.Debug("stdout: " + result.Stdout)
	return result, nil
}

func (r *ExecRunner) Whoami() (string, error) {
	result, err := r.Run("whoami")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

func formatCommand(name string, args []string) string {
	if len(args) == 0 {
		return name
	}
	return fmt.Sprintf("%s %s", name, strings.Join(args, " "))
}
