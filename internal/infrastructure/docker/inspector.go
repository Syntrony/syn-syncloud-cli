package docker

import (
	"bufio"
	"encoding/json"
	"strings"
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	domainresource "synctl/internal/domain/Resource"
	"synctl/internal/domain/docker"

	"github.com/google/uuid"
)

type Inspector struct {
	runner interfaces.CommandRunner
	logger interfaces.Logger
}

func NewInspector(runner interfaces.CommandRunner, logger interfaces.Logger) *Inspector {
	return &Inspector{
		runner: runner,
		logger: logger,
	}
}

func (i *Inspector) Snapshot() (*domain.Snapshot, error) {

	i.logger.Info("Collecting docker resources...")

	snapshot := &domain.Snapshot{}

	out, err := i.runner.Run(
		"docker",
		"ps",
		"-a",
		"--format",
		"{{json .}}",
	)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(out.Stdout))

	for scanner.Scan() {

		line := scanner.Text()

		var c docker.Container

		if err := json.Unmarshal([]byte(line), &c); err != nil {
			return nil, err
		}

		if strings.HasPrefix(c.Names, "k3d-") ||
			strings.HasPrefix(c.Names, "buildx_buildkit") ||
			strings.HasPrefix(c.Names, "devcontainer") {
			continue
		}

		resource := domain.Resource{
			Id:      c.ID,
			Name:    c.Names,
			Kind:    "docker.container",
			Runtime: "docker",
			NodeId:  uuid.NewString(),
			Spec: domainresource.Spec{
				Raw: map[string]interface{}{
					"image":   c.Image,
					"command": c.Command,
					"ports":   c.Ports,
					"labels":  c.Labels,
				},
			},
			Status: domainresource.Status{
				Raw: map[string]interface{}{
					"state": c.State,
				},
			},
			CreatedAt: c.Created,
			UpdatedAt: c.Created,
		}

		snapshot.Resources = append(snapshot.Resources, resource)

	}

	return snapshot, nil
}
