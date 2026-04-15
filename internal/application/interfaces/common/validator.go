package common

import "synctl/internal/domain"

type Validator interface {
	Validate(resource *domain.Resource) error
	ValidateSet(resource []*domain.Resource) error
}
