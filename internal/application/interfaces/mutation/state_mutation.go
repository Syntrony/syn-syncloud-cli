package mutation

import (
	"synctl/internal/domain"
	resource "synctl/internal/domain/Resource"
)

type StateMutation interface {
	Mutate(resources []*domain.Resource) ([]domain.Resource, error)
	MutateDns(incoming []resource.Record, existing []resource.Record) []resource.Record
}
