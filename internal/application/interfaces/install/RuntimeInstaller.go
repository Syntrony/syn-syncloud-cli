package install

import (
	"synctl/internal/domain"
)

type RuntimeInstaller interface {
	Install(mode domain.InstallMode) (*domain.Snapshot, error)
}
