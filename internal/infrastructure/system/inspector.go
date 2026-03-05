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

	node := domain.Node{
		Id:       uuid.NewString(),
		Ip:       ip,
		Hostname: hostname,
		Role:     "control-plane",
	}

	return []domain.Node{node}, nil

}
