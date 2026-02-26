package state

import (
	"synctl/internal/domain"
)

type Repository interface {
	Exists() (bool, error)
	Load() (*domain.State, error)
}
