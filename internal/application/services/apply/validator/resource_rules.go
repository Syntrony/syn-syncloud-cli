package validator

import (
	"fmt"

	"synctl/internal/domain"
)

func (v *ResourceValidator) validateBasicFields(resource *domain.Resource) error {
	if resource.Name == "" {
		return fmt.Errorf("resource name is required")
	}

	if resource.Kind == "" {
		return fmt.Errorf("resource kind is required")
	}

	if resource.Runtime == "" {
		return fmt.Errorf("runtime not resolved for kind: %s", resource.Kind)
	}

	return nil
}

func (v *ResourceValidator) validateRuntime(resource *domain.Resource) error {
	switch resource.Runtime {
	case "docker", "kubernetes":
		return nil
	default:
		return fmt.Errorf("unsupported runtime: %s", resource.Runtime)
	}
}

func (v *ResourceValidator) validateSpec(resource *domain.Resource) error {
	if resource.Spec == nil {
		return fmt.Errorf("spec cannot be empty for resource: %s", resource.Name)
	}

	switch {
	case resource.Kind == "docker.container":
		return v.validateDockerContainerSpec(resource.Spec)
	case resource.Kind == "docker.network":
		return v.validateDockerNetworkSpec(resource.Spec)
	case resource.Kind == "docker.image":
		return v.validateDockerImageSpec(resource.Spec)
	case resource.Kind == "k8s.raw":
		return v.validateK8sRawSpec(resource.Spec)
	}

	return nil
}

func (v *ResourceValidator) validateDockerContainerSpec(spec map[string]interface{}) error {
	if _, ok := spec["image"]; !ok {
		return fmt.Errorf("docker.container requires 'image' in spec")
	}
	return nil
}

func (v *ResourceValidator) validateDockerNetworkSpec(spec map[string]interface{}) error {
	return nil
}

func (v *ResourceValidator) validateDockerImageSpec(spec map[string]interface{}) error {
	return nil
}

func (v *ResourceValidator) validateK8sRawSpec(spec map[string]interface{}) error {
	if _, ok := spec["manifest"]; !ok {
		return fmt.Errorf("k8s.raw requires 'manifest' in spec")
	}
	return nil
}
