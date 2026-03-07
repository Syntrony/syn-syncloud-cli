package install

import (
	"synctl/internal/domain"
)

type RuntimeInstaller interface {
	Install() (*domain.Snapshot, error)
}
