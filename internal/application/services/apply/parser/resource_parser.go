package parser

import (
	"fmt"
	"synctl/internal/domain"
)

type ResourceParser struct {
	loader   *YamlLoader
	resolver *KindResolver
	builder  *ResourceBuilder
}

func NewResourceParser() *ResourceParser {
	return &ResourceParser{
		loader:   &YamlLoader{},
		resolver: NewKindResolver(),
		builder:  &ResourceBuilder{},
	}
}

func (p *ResourceParser) Parse(path string) ([]*domain.Resource, error) {
	dtos, err := p.loader.LoadAll(path)

	if err != nil {
		return nil, err
	}

	var resources []*domain.Resource

	for _, dto := range dtos {
		if dto.APIVersion == "" {
			return nil, fmt.Errorf("apiVersion required")
		}

		if dto.Metadata.Name == "" {
			return nil, fmt.Errorf("metadata.name required")
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
