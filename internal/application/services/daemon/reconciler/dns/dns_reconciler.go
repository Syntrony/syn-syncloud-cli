package dns

import (
	"sort"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	resource "synctl/internal/domain/Resource"
	dnsinfra "synctl/internal/infrastructure/dns"
	sysinfra "synctl/internal/infrastructure/system"
)

type DnsReconciler struct {
	inspector    *dnsinfra.Inspector
	installer    *dnsinfra.Installer
	sysInspector *sysinfra.Inspector
	logger       interfaces.Logger
}

func NewDnsReconciler(
	inspector *dnsinfra.Inspector,
	installer *dnsinfra.Installer,
	sysInspector *sysinfra.Inspector,
	logger interfaces.Logger,
) *DnsReconciler {
	return &DnsReconciler{
		inspector:    inspector,
		installer:    installer,
		sysInspector: sysInspector,
		logger:       logger,
	}
}

func (d *DnsReconciler) ReconcileDns(dns *domain.Dns) error {
	if dns == nil || len(dns.Records) == 0 {
		return nil
	}

	snap, err := d.inspector.Snapshot()
	if err != nil {
		d.logger.Error("DNS reconciler: failed to read current dnsmasq state")
		return err
	}

	if recordsEqual(snap.Records, dns.Records) {
		return nil
	}

	d.logger.Info("DNS reconciler: drift detected, applying desired records...")

	hosts, err := d.sysInspector.Inspect()
	if err != nil || len(hosts) == 0 {
		d.logger.Error("DNS reconciler: could not determine host IP")
		return err
	}

	return d.installer.ApplyRecords(dns.Records, hosts[0].Ip)
}

func recordsEqual(a, b []resource.Record) bool {
	if len(a) != len(b) {
		return false
	}
	sortRecords(a)
	sortRecords(b)
	for i := range a {
		if a[i].Name != b[i].Name || a[i].Server != b[i].Server {
			return false
		}
	}
	return true
}

func sortRecords(records []resource.Record) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})
}
