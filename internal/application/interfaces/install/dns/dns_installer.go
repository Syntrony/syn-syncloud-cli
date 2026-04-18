package dns

import resource "synctl/internal/domain/resource"

type DnsInstaller interface {
	InstallDns() error
	ConfigureDns() error
	SetupDnsFile(records []resource.Record, generalIp string) string
	// WithConfDir overrides the directory where syncloud writes its dnsmasq config file.
	// When not set, defaults to /etc/dnsmasq.d.
	WithConfDir(dir string)
}
