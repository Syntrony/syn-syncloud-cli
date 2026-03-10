package dns

import (
	"bufio"
	"os"
	"path/filepath"
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

	files, err := filepath.Glob("/etc/dnsmasq.d/*.conf")

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// 2. Abrimos cada archivo individualmente
		f, err := os.Open(file)
		if err != nil {
			i.logger.Info("Could not read file: " + file)
			continue
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Tu lógica de parsing se mantiene igual
			if strings.HasPrefix(line, "address=/") {
				trimmed := strings.TrimPrefix(line, "address=/")
				parts := strings.Split(trimmed, "/")

				if len(parts) == 2 {
					snapshot.Records = append(snapshot.Records, resource.Record{
						Name:   parts[0],
						Server: parts[1],
					})
				}
			}
		}
		f.Close()
	}

	return snapshot, nil

}
