package install

import (
	"synctl/internal/domain"
)

type Inspector interface {
	Snapshot() (*domain.Snapshot, error)
}
