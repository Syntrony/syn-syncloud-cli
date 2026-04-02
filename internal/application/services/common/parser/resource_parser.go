package parser

import (
	"errors"
	"fmt"
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
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	dtos, err := p.loader.LoadAll(path)

	if err != nil {
		return nil, fmt.Errorf("failed to load YAML file: %w", err)
	}

	var resources []*domain.Resource

	for _, dto := range dtos {
		if dto.APIVersion == "" {
			return nil, fmt.Errorf("invalid resource: APIVersion is required for kind %s", dto.Kind)
		}

		if dto.Metadata.Name == "" {
			return nil, fmt.Errorf("invalid resource: metadata.name is required for kind %s", dto.Kind)
		}

		mapping, err := p.resolver.Resolve(dto.Kind)

		if err != nil {
			return nil, fmt.Errorf("unsupported resource kind %s: %w", dto.Kind, err)
		}

		resource := p.builder.Build(&dto, mapping)
		resources = append(resources, resource)
	}

	return resources, nil
}
