package executor

import (
	"fmt"
	"os/exec"
	"syscall"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"

	"golang.org/x/term"
)

type SudoRunner struct {
	logger interfaces.Logger
}

func NewSudoRunner(logger interfaces.Logger) *SudoRunner {
	return &SudoRunner{logger: logger}
}

func (r *SudoRunner) Run(cmd string, args ...string) (*dto.CommandResult, error) {

	fmt.Print("[sudo] password: ")

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return nil, err
	}

	password := string(bytePassword)

	fullArgs := append([]string{"-S", cmd}, args...)

	command := exec.Command("sudo", fullArgs...)

	stdin, _ := command.StdinPipe()

	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, password)
	}()

	output, err := command.CombinedOutput()

	return &dto.CommandResult{
		Stdout: string(output),
	}, err
}
