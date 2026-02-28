package interfaces

import (
	"synctl/internal/domain"
)

type Inspector interface {
	Snapshot() (*domain.Snapshot, error)
}
