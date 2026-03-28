package parser

import (
	"synctl/internal/domain"
)

type ResourceParser interface {
	Parse(file string) ([]*domain.Resource, error)
}

type resourceParser struct {
	loader   *YamlLoader
	resolver *KindResolver
	builder  *ResourceBuilder
}

func NewResourceParser() ResourceParser {
	return &resourceParser{
		loader:   &YamlLoader{},
		resolver: NewKindResolver(),
		builder:  &ResourceBuilder{},
	}
}

func (p *resourceParser) Parse(path string) ([]*domain.Resource, error) {
	dtos, err := p.loader.LoadAll(path)

	if err != nil {
		return nil, err
	}

	var resources []*domain.Resource

	for _, dto := range dtos {
		if dto.APIVersion == "" {
			return nil, err
		}

		if dto.Metadata.Name == "" {
			return nil, err
		}

		mapping, err := p.resolver.Resolve(dto.Kind)

		if err != nil {
			return nil, err
		}

		resource := p.builder.Build(&dto, mapping)
		resources = append(resources, resource)
	}

	return resources, nil
}
