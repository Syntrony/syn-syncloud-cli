package dns

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	resource "synctl/internal/domain/resource"
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

	// Use subset check: all desired records must exist in actual.
	// External records managed outside syncloud are intentionally ignored.
	if desiredSubsetOfActual(dns.Records, snap.Records) {
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

// desiredSubsetOfActual returns true when every record in desired exists in actual
// with the same server value. Extra records in actual (e.g. managed by external
// tools in /etc/dnsmasq.conf) are intentionally ignored to avoid phantom drift.
func desiredSubsetOfActual(desired, actual []resource.Record) bool {
	index := make(map[string]string, len(actual))
	for _, r := range actual {
		index[r.Name] = r.Server
	}
	for _, r := range desired {
		if server, ok := index[r.Name]; !ok || server != r.Server {
			return false
		}
	}
	return true
}
