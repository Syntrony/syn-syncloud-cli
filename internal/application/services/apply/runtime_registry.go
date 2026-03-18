package apply

import (
	"synctl/internal/application/interfaces"
	iApply "synctl/internal/application/interfaces/apply"
	"synctl/internal/domain"
)

type RuntimeRegistry struct {
	runtimes map[string]iApply.RuntimeReconciler
}

func NewRuntimeRegistry() *RuntimeRegistry {
	return &RuntimeRegistry{
		runtimes: make(map[string]iApply.RuntimeReconciler),
	}
}

func (r *RuntimeRegistry) Register(reconciler iApply.RuntimeReconciler) {
	r.runtimes[reconciler.Runtime()] = reconciler
}

func (r *RuntimeRegistry) Get(runtimeName string) (iApply.RuntimeReconciler, bool) {
	reconciler, ok := r.runtimes[runtimeName]
	return reconciler, ok
}

func (r *RuntimeRegistry) All() []iApply.RuntimeReconciler {
	result := make([]iApply.RuntimeReconciler, 0, len(r.runtimes))
	for _, recon := range r.runtimes {
		result = append(result, recon)
	}
	return result
}

func (r *RuntimeRegistry) FilterForRuntime(resources []*domain.Resource) map[string][]*domain.Resource {
	result := make(map[string][]*domain.Resource)
	for _, res := range resources {
		result[res.Runtime] = append(result[res.Runtime], res)
	}
	return result
}

func DefaultRuntimeRegistry(runner interfaces.CommandRunner, logger interfaces.Logger) *RuntimeRegistry {
	registry := NewRuntimeRegistry()

	// Note: Callers should register reconcilers from their packages
	// This is a helper to create an empty registry
	return registry
}
