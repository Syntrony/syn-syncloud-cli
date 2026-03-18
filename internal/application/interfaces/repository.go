package interfaces

import (
	"synctl/internal/domain"
)

type Repository interface {
	Exists() (bool, error)
	Load() (*domain.State, error)
	Save(st *domain.State) error
	Upsert(resources []*domain.Resource) (*domain.Snapshot, error)
}
