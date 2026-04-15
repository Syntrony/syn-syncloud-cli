package validator

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
)

type Validator interface {
	Validate(resource *domain.Resource) error
	ValidateSet(resource []*domain.Resource) error
}

type resourceValidator struct {
	logger interfaces.Logger
}

func NewValidator() Validator {
	return &resourceValidator{}
}

func (v *resourceValidator) Validate(resource *domain.Resource) error {
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

func (v *resourceValidator) ValidateSet(resources []*domain.Resource) error {
	if err := v.validateDuplicateNames(resources); err != nil {
		return err
	}

	if err := v.validateNamespaceRules(resources); err != nil {
		return err
	}

	return nil
}
