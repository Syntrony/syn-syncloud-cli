package common

import "synctl/internal/domain"

type ResourceParser interface {
	Parse(file string) ([]*domain.Resource, error)
}
