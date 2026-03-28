package mutation

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
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
