package common

import "synctl/internal/domain"

type DnsReconciler interface {
	ReconcileDns(dns *domain.Dns) error
}
