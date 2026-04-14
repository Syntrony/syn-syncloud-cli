package executor

import (
	"os"
	"os/exec"
	"time"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"
)

type SudoRunner struct {
	runner       interfaces.CommandRunner
	lastAuthTime time.Time
	authTimeout  time.Duration
}

func NewSudoRunner(runner interfaces.CommandRunner) *SudoRunner {
	return &SudoRunner{
		runner:      runner,
		authTimeout: 10 * time.Minute,
	}
}

func (r *SudoRunner) EnsureAuth() error {
	now := time.Now()
	if now.Sub(r.lastAuthTime) < r.authTimeout {
		return nil
	}

	// Use exec directly with stdin so user can enter password interactively if needed
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	r.lastAuthTime = now
	return nil
}

func (r *SudoRunner) Run(cmd string, args ...string) (*dto.CommandResult, error) {
	if err := r.EnsureAuth(); err != nil {
		return nil, err
	}

	fullArgs := append([]string{cmd}, args...)

	return r.runner.Run("sudo", fullArgs...)
}

func (r *SudoRunner) Whoami() (string, error) {
	return r.runner.Whoami()
}
