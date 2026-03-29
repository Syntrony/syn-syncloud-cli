package services

import (
	interfaces "synctl/internal/application/interfaces"
	"synctl/internal/domain"
	dto "synctl/internal/domain/dto"
	filters "synctl/internal/domain/filters"
)

type GetResourceService struct {
	Repo interfaces.Repository
}

func (s *GetResourceService) Execute(filter filters.GetResourceFilter) ([]dto.ResourceDto, error) {
	resources, err := s.GetResourceBy(filter)

	if err != nil {
		return nil, err
	}

	var result []dto.ResourceDto

	for _, r := range resources {

		result = append(result, dto.ResourceDto{
			Id:        r.Id,
			Name:      r.Name,
			Runtime:   r.Runtime,
			CreatedAt: r.CreatedAt,
		})
	}

	return result, nil

}

func (s *GetResourceService) GetResourceBy(filter filters.GetResourceFilter) ([]*domain.Resource, error) {
	exists, err := s.Repo.Exists()
	var resources []*domain.Resource

	if err != nil {
		return nil, err
	}

	if !exists {
		return resources, nil
	}

	st, err := s.Repo.Load()

	if err != nil {
		return nil, err
	}

	for _, r := range st.Resources {
		if (filter.Name == "" || r.Name == filter.Name) &&
			(filter.Kind == "" || r.Kind == filter.Kind) &&
			(filter.Id == "" || r.Id == filter.Id) &&
			(filter.Runtime == "" || r.Runtime == filter.Runtime) &&
			(filter.NodeId == "" || r.NodeId == filter.NodeId) {
			resources = append(resources, &r)
		}
	}

	return resources, nil
}
