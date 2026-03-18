package validator

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
)

type ResourceValidator struct {
	logger interfaces.Logger
}

func NewValidator() *ResourceValidator {
	return &ResourceValidator{}
}

func (v *ResourceValidator) Validate(resource *domain.Resource) error {
	if err := v.validateBasicFields(resource); err != nil {
		return err
	}

	if err := v.validateRuntime(resource); err != nil {
		return err
	}

	if err := v.validateSpec(resource); err != nil {
		return err
	}

	return nil
}

func (v *ResourceValidator) ValidateSet(resources []*domain.Resource) error {
	if err := v.validateDuplicateNames(resources); err != nil {
		return err
	}

	if err := v.validateNamespaceRules(resources); err != nil {
		return err
	}

	return nil
}
