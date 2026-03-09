package dns

import (
	"bufio"
	"strings"
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	resource "synctl/internal/domain/Resource"
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
	i.logger.Info("Collecting dns resources...")
	snapshot := &domain.Snapshot{}

	out, err := i.runner.Run(
		"cat",
		"/etc/dnsmasq.d/syncloud.conf",
	)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(out.Stdout))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !(line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "server=")) && strings.HasPrefix(line, "address=/") {
			// rm prefix
			trimmed := strings.TrimPrefix(line, "address=/")

			// split "google.com/1.2.3.4"
			parts := strings.Split(trimmed, "/")

			if len(parts) == 2 {
				snapshot.Records = append(snapshot.Records, resource.Record{
					Name:   parts[0],
					Server: parts[1],
				})
			}
		}
	}

	return snapshot, nil

}
