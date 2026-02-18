package system

import (
	"bytes"
	"errors"
	"os/exec"
)

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	return out.String(), nil
}
