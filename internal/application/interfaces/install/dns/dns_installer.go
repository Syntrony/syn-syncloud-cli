package dns

import resource "synctl/internal/domain/resource"

type DnsInstaller interface {
	InstallDns() error
	ConfigureDns() error
	SetupDnsFile(records []resource.Record, generalIp string) string
}
