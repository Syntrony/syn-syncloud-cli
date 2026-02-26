package services

import (
	dto "synctl/internal/domain/Dto"
	filters "synctl/internal/domain/filters"
	repository "synctl/internal/infrastructure/state"
)

type GetResourceService struct {
	Repo repository.Repository
}

func (s *GetResourceService) Execute(filter filters.GetResourceFilter) ([]dto.ResourceDto, error) {
	exists, err := s.Repo.Exists()

	if err != nil {
		return nil, err
	}

	if !exists {
		return []dto.ResourceDto{}, nil
	}

	st, err := s.Repo.Load()
	if err != nil {
		return nil, err
	}

	var resources []dto.ResourceDto

	for _, r := range st.Resources {
		if filter.Name != "" && r.Name != filter.Name {
			continue
		}
		if filter.Kind != "" && r.Kind != filter.Kind {
			continue
		}
		if filter.Id != "" && r.Id != filter.Id {
			continue
		}
		if filter.Runtime != "" && r.Runtime != filter.Runtime {
			continue
		}
		if filter.NodeId != "" && r.NodeId != filter.NodeId {
			continue
		}

		resources = append(resources, dto.ResourceDto{
			Id:        r.Id,
			Name:      r.Name,
			Runtime:   r.Runtime,
			CreatedAt: r.CreatedAt,
		})
	}

	return resources, nil

}
