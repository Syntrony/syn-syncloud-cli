package mutation

import (
	"synctl/internal/domain"
)

type StateMutation interface {
	Mutate(resources []*domain.Resource) ([]domain.Resource, error)
}
