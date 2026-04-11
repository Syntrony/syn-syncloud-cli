package describe

import (
	"fmt"
	"synctl/internal/application/interfaces"
	"synctl/internal/application/interfaces/common"
	services "synctl/internal/application/services/get"
	"synctl/internal/domain"
	"synctl/internal/domain/dto"
	"synctl/internal/domain/filters"
)

type DescribeService struct {
	Repo     interfaces.Repository
	Runtimes []common.RuntimeReconciler
	GetSvc   services.GetResourceService
}

func NewDescribeService(repo interfaces.Repository, runtimes []common.RuntimeReconciler, getSvc services.GetResourceService) *DescribeService {
	return &DescribeService{
		Repo:     repo,
		Runtimes: runtimes,
		GetSvc:   getSvc,
	}
}

func (s *DescribeService) Execute(id, name, kind, runtime string) (dto.ResourceDetailDto, error) {
	filter := filters.GetResourceFilter{
		Name:    name,
		Kind:    kind,
		Id:      id,
		Runtime: runtime,
	}

	resources, err := s.GetSvc.GetResourceBy(filter)

	if err != nil {
		return dto.ResourceDetailDto{}, err
	}

	if len(resources) == 0 {
		return s.FetchFromRuntime(name, kind, runtime)
	}

	resource := resources[0]
	liveSpec, err := s.ObserveRuntime(resource)

	if err != nil {
		return s.BuildDtoFromState(resource, nil), nil
	}

	return s.BuildDtoFromState(resource, liveSpec), nil

}

func (s *DescribeService) ObserveRuntime(resource *domain.Resource) (map[string]interface{}, error) {
	for _, runtime := range s.Runtimes {
		if runtime.Runtime() == resource.Runtime {
			state, err := runtime.Observe(resource)

			if err != nil {
				return nil, err
			}

			if state.Exists {
				return state.Spec, nil
			}
		}
	}
	return nil, nil
}

func (s *DescribeService) FetchFromRuntime(name, kind, runtime string) (dto.ResourceDetailDto, error) {
	res := &domain.Resource{
		Name:    name,
		Kind:    kind,
		Runtime: runtime,
	}

	for _, runtime := range s.Runtimes {
		state, err := runtime.Observe(res)
		if err != nil {
			continue
		}

		if state.Exists {
			return dto.ResourceDetailDto{
				Name:     name,
				Kind:     kind,
				Runtime:  runtime.Runtime(),
				LiveSpec: state.Spec,
			}, nil
		}
	}

	return dto.ResourceDetailDto{}, fmt.Errorf("resource not found: %s/%s", runtime, name)
}

func (s *DescribeService) BuildDtoFromState(resource *domain.Resource, liveSpec map[string]interface{}) dto.ResourceDetailDto {
	return dto.ResourceDetailDto{
		Id:        resource.Id,
		Name:      resource.Name,
		Kind:      resource.Kind,
		Runtime:   resource.Runtime,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
		Status:    resource.Status,
		Spec:      resource.Spec,
		LiveSpec:  liveSpec,
	}
}
