package dns

import resource "synctl/internal/domain/Resource"

type DnsInstaller interface {
	InstallDns() error
	ConfigureDns() error
	SetupDnsFile(records []resource.Record, generalIp string) string
}
