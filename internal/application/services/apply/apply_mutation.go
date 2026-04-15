package mutation

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	resource "synctl/internal/domain/resource"
)

type ApplyMutation struct {
	repo interfaces.Repository
}

func NewApplyMutation(repo interfaces.Repository) *ApplyMutation {
	return &ApplyMutation{
		repo: repo,
	}
}

func (m *ApplyMutation) Mutate(resources []*domain.Resource) ([]domain.Resource, error) {
	return m.repo.Upsert(resources)
}

func (m *ApplyMutation) MutateDns(incoming []resource.Record, existing []resource.Record) []resource.Record {
	merged := make(map[string]resource.Record)
	for _, r := range existing {
		merged[r.Name] = r
	}
	for _, r := range incoming {
		merged[r.Name] = r
	}
	result := make([]resource.Record, 0, len(merged))
	for _, r := range merged {
		result = append(result, r)
	}
	return result
}
