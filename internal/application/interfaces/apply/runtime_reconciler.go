package apply

import "synctl/internal/domain"

type RuntimeReconciler interface {
	Runtime() string
	Observe(resource *domain.Resource) (*domain.RuntimeResourceState, error)
	Reconcile(context domain.ReconcileContext) error
}
