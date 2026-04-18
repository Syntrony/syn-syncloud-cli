package dns

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	resource "synctl/internal/domain/resource"
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

	files := i.discoverConfFiles()

	seen := make(map[string]struct{})
	for _, file := range files {
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}

		records, err := i.parseAddressRecords(file)
		if err != nil {
			i.logger.Info("Could not read dns config file: " + file)
			continue
		}
		snapshot.Records = append(snapshot.Records, records...)
	}

	return snapshot, nil
}

// discoverConfFiles returns all dnsmasq config file paths by:
//  1. Reading /etc/dnsmasq.conf (main config) and following conf-dir= / conf-file= directives.
//  2. Globbing the standard drop-in directory /etc/dnsmasq.d/*.conf.
func (i *Inspector) discoverConfFiles() []string {
	var files []string

	mainConf := "/etc/dnsmasq.conf"
	if _, err := os.Stat(mainConf); err == nil {
		files = append(files, mainConf)
		files = append(files, i.resolveDirectives(mainConf)...)
	}

	dropins, _ := filepath.Glob("/etc/dnsmasq.d/*.conf")
	files = append(files, dropins...)

	return files
}

// resolveDirectives reads a dnsmasq config file and returns additional paths
// referenced by conf-dir= and conf-file= directives.
func (i *Inspector) resolveDirectives(configFile string) []string {
	var extra []string

	f, err := os.Open(configFile)
	if err != nil {
		return extra
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(line, "conf-dir="):
			dir := strings.TrimPrefix(line, "conf-dir=")
			// conf-dir= accepts optional suffix filters: conf-dir=/path,.conf
			parts := strings.SplitN(dir, ",", 2)
			dirPath := strings.TrimSpace(parts[0])
			pattern := dirPath + "/*.conf"
			if len(parts) > 1 {
				ext := strings.TrimSpace(parts[1])
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				pattern = dirPath + "/*" + ext
			}
			matches, _ := filepath.Glob(pattern)
			extra = append(extra, matches...)

		case strings.HasPrefix(line, "conf-file="):
			filePath := strings.TrimSpace(strings.TrimPrefix(line, "conf-file="))
			if _, err := os.Stat(filePath); err == nil {
				extra = append(extra, filePath)
			}
		}
	}

	return extra
}

// parseAddressRecords extracts all address=/domain/ip entries from a single config file.
func (i *Inspector) parseAddressRecords(file string) ([]resource.Record, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []resource.Record
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "address=/") {
			continue
		}
		trimmed := strings.TrimPrefix(line, "address=/")
		parts := strings.Split(trimmed, "/")
		if len(parts) == 2 {
			records = append(records, resource.Record{
				Name:   parts[0],
				Server: parts[1],
			})
		}
	}

	return records, nil
}
