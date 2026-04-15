package diff

import "synctl/internal/domain"

func Compare(desired *domain.Resource, actual *domain.RuntimeResourceState) *domain.Action {
	if !actual.Exists {
		return &domain.Action{
			Type:     domain.ActionCreate,
			Resource: desired,
		}
	}

	return &domain.Action{
		Type:     domain.ActionNoop,
		Resource: desired,
	}
}

func CompareAndUpdate(desired *domain.Resource, actual *domain.RuntimeResourceState) *domain.Action {
	if !actual.Exists {
		return &domain.Action{
			Type:     domain.ActionCreate,
			Resource: desired,
		}
	}

	return &domain.Action{
		Type:     domain.ActionNoop,
		Resource: desired,
	}
}
