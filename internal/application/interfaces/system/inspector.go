package system

import "synctl/internal/domain"

type Inspector interface {
	Inspect() ([]domain.Node, error)
}
