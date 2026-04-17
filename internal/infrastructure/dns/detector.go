package dns

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"
)

type Detector struct {
	runner interfaces.CommandRunner
	logger interfaces.Logger
}

func NewDetector(runner interfaces.CommandRunner, logger interfaces.Logger) *Detector {
	return &Detector{
		runner: runner,
		logger: logger,
	}
}

func (d *Detector) Detect() (*dto.Status, error) {
	d.logger.Info("Detecting dns environment...")
	status := &dto.Status{}

	_, err := d.runner.Run("dnsmasq", "--version")

	if err != nil {
		status.DnsInstalled = false
		return status, nil
	}

	status.DnsInstalled = true

	// Use sh to expand the glob — exec.Command does not invoke a shell,
	// so passing "*.conf" directly would look for a file literally named "*.conf".
	_, err = d.runner.Run("sh", "-c", "ls /etc/dnsmasq.d/*.conf")

	if err != nil {
		status.DnsConfigured = false
		return status, nil
	}

	status.DnsConfigured = true

	return status, nil
}
