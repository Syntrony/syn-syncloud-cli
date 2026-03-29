package reconciler

import (
	"fmt"
	"synctl/internal/application/interfaces/common"
	"synctl/internal/application/services/common/diff"
	"synctl/internal/domain"
)

type ReconcileService struct {
	runtimes []common.RuntimeReconciler
}

func NewReconcileService(runtimes []common.RuntimeReconciler) *ReconcileService {
	return &ReconcileService{
		runtimes: runtimes,
	}
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

		if err := runtime.Reconcile(subCtx); err != nil {
			return err
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
