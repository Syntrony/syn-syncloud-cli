package mutation

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
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
