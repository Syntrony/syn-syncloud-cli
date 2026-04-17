package reconciler

import (
	"fmt"
	"synctl/internal/application/interfaces"
	"synctl/internal/application/interfaces/common"
	"synctl/internal/application/services/common/diff"
	"synctl/internal/domain"
)

type ReconcileService struct {
	runtimes      []common.RuntimeReconciler
	dnsReconciler common.DnsReconciler
	dryRun        bool
	logger        interfaces.Logger
}

func NewReconcileService(runtimes []common.RuntimeReconciler) *ReconcileService {
	return &ReconcileService{
		runtimes: runtimes,
	}
}

func (r *ReconcileService) WithDnsReconciler(dns common.DnsReconciler) *ReconcileService {
	r.dnsReconciler = dns
	return r
}

func (r *ReconcileService) WithDryRun(logger interfaces.Logger) *ReconcileService {
	r.dryRun = true
	r.logger = logger
	return r
}

func (r *ReconcileService) Reconcile(state *domain.State) error {
	ctx := domain.ReconcileContext{
		Resources: toResourcePointers(state.Resources),
		State:     state,
	}

	for _, runtime := range r.runtimes {
		var runtimeResources []*domain.Resource

		for _, resource := range ctx.Resources {
			if resource.Runtime == runtime.Runtime() {
				runtimeResources = append(runtimeResources, resource)
			}
		}

		if len(runtimeResources) == 0 {
			continue
		}

		var planned []*domain.Action

		for _, resource := range runtimeResources {
			actual, err := runtime.Observe(resource)
			if err != nil {
				return fmt.Errorf("observing resource %s: %w", resource.Name, err)
			}

			action := diff.Compare(resource, actual)
			planned = append(planned, action)
		}

		subCtx := domain.ReconcileContext{
			Resources: runtimeResources,
			State:     ctx.State,
			Actions:   planned,
		}

		if r.dryRun {
			for _, action := range planned {
				r.logger.Info(fmt.Sprintf("[dry-run] %s → %s %s", runtime.Runtime(), action.Type, action.Resource.Name))
			}
			continue
		}

		if err := runtime.Reconcile(subCtx); err != nil {
			return err
		}
	}

	if r.dnsReconciler != nil && state.Dns != nil {
		if r.dryRun {
			r.logger.Info("[dry-run] dns reconciliation skipped — no changes applied")
		} else if err := r.dnsReconciler.ReconcileDns(state.Dns); err != nil {
			return fmt.Errorf("dns reconciliation: %w", err)
		}
	}

	return nil
}

func toResourcePointers(resources []domain.Resource) []*domain.Resource {
	result := make([]*domain.Resource, 0, len(resources))

	for i := range resources {
		result = append(result, &resources[i])
	}

	return result
}
