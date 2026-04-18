package domain

import (
	"fmt"
	resource "synctl/internal/domain/resource"
)

// SourceManaged indicates a resource declared explicitly via synctl apply.
// Only managed resources are reconciled (created/updated) by the reconciler.
const SourceManaged = "managed"

// SourceObserved indicates a resource imported/discovered by an inspector.
// Observed resources are stored in state.json for visibility but are NOT reconciled.
const SourceObserved = "observed"

type Resource struct {
	Id        string                 `json:"id"`
	Name      string                 `json:"name"`
	Kind      string                 `json:"kind"`
	Runtime   string                 `json:"runtime"`
	NodeId    string                 `json:"nodeId"`
	Source    string                 `json:"source"`
	Spec      map[string]interface{} `json:"spec"`
	Status    map[string]interface{} `json:"status"`
	Ownership resource.Ownership     `json:"ownership"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

func (r *Resource) Key() string {
	return fmt.Sprintf("%s/%s/%s",
		r.Runtime,
		r.Kind,
		r.Name,
	)
}
