package system

import (
	"os"
	"strings"
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"

	"github.com/google/uuid"
)

type Inspector struct {
	runner interfaces.CommandRunner
}

func NewInspector(runner interfaces.CommandRunner) *Inspector {
	return &Inspector{
		runner: runner,
	}
}

func (i *Inspector) Inspect() ([]domain.Node, error) {
	hostname, _ := os.Hostname()

	out, err := i.runner.Run("hostname", "-I")
	if err != nil {
		return nil, err
	}

	ip := strings.Fields(out.Stdout)[0]

	nodeID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(hostname)).String()

	osInfo := i.readOsRelease()
	arch := i.readArch()

	node := domain.Node{
		Id:       nodeID,
		Ip:       ip,
		Hostname: hostname,
		Role:     "control-plane",
		Os:       osInfo,
		Arch:     arch,
	}

	return []domain.Node{node}, nil
}

func (i *Inspector) readOsRelease() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			val := strings.TrimPrefix(line, "PRETTY_NAME=")
			return strings.Trim(val, `"`)
		}
	}
	return "unknown"
}

func (i *Inspector) readArch() string {
	out, err := i.runner.Run("uname", "-m")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(out.Stdout)
}
