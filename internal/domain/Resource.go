package domain

import (
	resource "synctl/internal/domain/Resource"
)

type Resource struct {
	Id        string             `json:"id"`
	Name      string             `json:"name"`
	Kind      string             `json:"kind"`
	Runtime   string             `json:"runtime"`
	NodeId    string             `json:"nodeId"`
	Spec      resource.Spec      `json:"spec"`
	Status    resource.Status    `json:"status"`
	Ownership resource.Ownership `json:"ownership"`
	CreatedAt string             `json:"createdAt"`
	UpdatedAt string             `json:"updatedAt"`
}
