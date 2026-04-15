package validator

import (
	"fmt"
	"synctl/internal/domain"
)

func (v *resourceValidator) validateDuplicateNames(resources []*domain.Resource) error {
	seen := make(map[string]bool)

	for _, resource := range resources {
		key := resource.Kind + ":" + resource.Name

		if seen[key] {
			errMsg := "duplicate resource: " + key
			v.logger.Error(fmt.Sprintf(errMsg))
			return fmt.Errorf(errMsg)
		}

		seen[key] = true
	}

	return nil
}

func (v *resourceValidator) validateNamespaceRules(resources []*domain.Resource) error {

	namespaces := map[string]bool{}

	for _, r := range resources {
		if r.Kind == "k8s.namespace" {
			namespaces[r.Name] = true
		}
	}

	for _, r := range resources {

		ns := r.Name
		if ns == "" {
			continue
		}

		if !namespaces[ns] {
			continue
		}
	}

	return nil
}
