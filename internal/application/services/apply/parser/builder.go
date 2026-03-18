package parser

import (
	"synctl/internal/domain"
	"synctl/internal/domain/parser"

	"github.com/google/uuid"
)

type ResourceBuilder struct{}

func (b *ResourceBuilder) Build(
	dto *parser.ResourceYAML,
	mapping KindMapping,
) *domain.Resource {
	return &domain.Resource{
		Id:      uuid.NewString(),
		Name:    dto.Metadata.Name,
		Runtime: mapping.Runtime,
		Kind:    mapping.Kind,
		Spec:    dto.Spec,
	}
}
