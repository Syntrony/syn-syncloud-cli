package mutation

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	resource "synctl/internal/domain/Resource"
)

type DeleteMutation struct {
	repo interfaces.Repository
}

func NewDeleteMutation(repo interfaces.Repository) *DeleteMutation {
	return &DeleteMutation{
		repo: repo,
	}
}

func (m *DeleteMutation) Mutate(resources []*domain.Resource) ([]domain.Resource, error) {
	return m.repo.Remove(resources)
}

func (m *DeleteMutation) MutateDns(incoming []resource.Record, existing []resource.Record) []resource.Record {
	remove := make(map[string]struct{})
	for _, r := range incoming {
		remove[r.Name] = struct{}{}
	}
	var result []resource.Record
	for _, r := range existing {
		if _, ok := remove[r.Name]; !ok {
			result = append(result, r)
		}
	}
	return result
}
