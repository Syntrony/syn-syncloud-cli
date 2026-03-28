package parser

import (
	"fmt"
)

type KindMapping struct {
	Runtime      string
	Kind         string
	Mode         string
	RequiredSpec []string
}

type KindResolver struct {
	registry map[string]KindMapping
}

func NewKindResolver() *KindResolver {
	return &KindResolver{
		registry: map[string]KindMapping{
			"Deployment": {
				Runtime:      "kubernetes",
				Kind:         "k8s.deployment",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Namespace": {
				Runtime:      "kubernetes",
				Kind:         "k8s.namespace",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Service": {
				Runtime:      "kubernetes",
				Kind:         "k8s.service",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Ingress": {
				Runtime:      "kubernetes",
				Kind:         "k8s.ingress",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"ServiceAccount": {
				Runtime:      "kubernetes",
				Kind:         "k8s.serviceaccount",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Role": {
				Runtime:      "kubernetes",
				Kind:         "k8s.role",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"RoleBinding": {
				Runtime:      "kubernetes",
				Kind:         "k8s.rolebinding",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"ClusterRole": {
				Runtime:      "kubernetes",
				Kind:         "k8s.clusterrole",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Secret": {
				Runtime:      "kubernetes",
				Kind:         "k8s.secret",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"ConfigMap": {
				Runtime:      "kubernetes",
				Kind:         "k8s.configmap",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Endpoints": {
				Runtime:      "kubernetes",
				Kind:         "k8s.endpoints",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Pod": {
				Runtime:      "kubernetes",
				Kind:         "k8s.pod",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"K8sResource": {
				Runtime:      "kubernetes",
				Kind:         "k8s.raw",
				Mode:         "raw",
				RequiredSpec: []string{"manifest"},
			},
			"Container": {
				Runtime:      "docker",
				Kind:         "docker.container",
				Mode:         "abstract",
				RequiredSpec: []string{"image"},
			},
			"Network": {
				Runtime:      "docker",
				Kind:         "docker.network",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
			"Image": {
				Runtime:      "docker",
				Kind:         "docker.image",
				Mode:         "abstract",
				RequiredSpec: []string{},
			},
		},
	}
}

func (r *KindResolver) Resolve(kind string) (KindMapping, error) {
	mapping, ok := r.registry[kind]

	if !ok {
		return KindMapping{}, fmt.Errorf("unsupported kind: %s", kind)
	}

	return mapping, nil
}

func (r *KindResolver) IsSupported(kind string) bool {
	_, ok := r.registry[kind]
	return ok
}

func (r *KindResolver) SupportedKinds() []string {
	kinds := make([]string, 0, len(r.registry))
	for kind := range r.registry {
		kinds = append(kinds, kind)
	}
	return kinds
}

func (r *KindResolver) RuntimeForKind(kind string) (string, bool) {
	mapping, ok := r.registry[kind]
	return mapping.Runtime, ok
}
