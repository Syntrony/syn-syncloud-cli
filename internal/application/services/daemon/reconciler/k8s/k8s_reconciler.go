package k8s

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
)

var supportedKinds = []string{
	"namespace",
	"deployment",
	"service",
	"ingress",
	"serviceaccount",
	"role",
	"rolebinding",
	"clusterrole",
	"secret",
	"configmap",
	"endpoints",
	"pod",
	"raw",
}

var clusterScopedKinds = map[string]bool{
	"namespace":   true,
	"clusterrole": true,
}

var apiVersionMap = map[string]string{
	"deployment":     "apps/v1",
	"service":        "v1",
	"namespace":      "v1",
	"ingress":        "networking.k8s.io/v1",
	"serviceaccount": "v1",
	"role":           "rbac.authorization.k8s.io/v1",
	"rolebinding":    "rbac.authorization.k8s.io/v1",
	"clusterrole":    "rbac.authorization.k8s.io/v1",
	"secret":         "v1",
	"configmap":      "v1",
	"endpoints":      "v1",
	"pod":            "v1",
}

type K8sReconciler struct {
	runner interfaces.CommandRunner
	logger interfaces.Logger
}

func NewK8sReconciler(runner interfaces.CommandRunner, logger interfaces.Logger) *K8sReconciler {
	return &K8sReconciler{
		runner: runner,
		logger: logger,
	}
}

func (k *K8sReconciler) Runtime() string {
	return "kubernetes"
}

func (k *K8sReconciler) Observe(resource *domain.Resource) (*domain.RuntimeResourceState, error) {
	kind := extractKindName(resource.Kind)
	namespace := extractNamespace(resource.Spec)

	var cmd []string
	if namespace != "" {
		cmd = []string{"get", kind, resource.Name, "-n", namespace, "-o", "json"}
	} else {
		cmd = []string{"get", kind, resource.Name, "-o", "json"}
	}

	out, err := k.runner.Run("kubectl", cmd...)
	if err != nil || out.Stdout == "" {
		return &domain.RuntimeResourceState{Exists: false}, nil
	}

	return &domain.RuntimeResourceState{
		Exists: true,
		Spec:   map[string]interface{}{"raw": out.Stdout},
	}, nil
}

func (k *K8sReconciler) Reconcile(ctx domain.ReconcileContext) error {
	for _, action := range ctx.Actions {
		res := action.Resource
		kind := extractKindName(res.Kind)

		if !isSupportedKind(kind) {
			k.logger.Info(fmt.Sprintf("Kind %s not supported yet, skipping", kind))
			continue
		}

		switch action.Type {
		case domain.ActionCreate:
			if err := k.applyResource(res); err != nil {
				return fmt.Errorf("failed to apply %s %s: %w", kind, res.Name, err)
			}
		case domain.ActionUpdate:
			if err := k.applyResource(res); err != nil {
				return fmt.Errorf("failed to update %s %s: %w", kind, res.Name, err)
			}
		case domain.ActionDelete:
			if err := k.deleteResource(res); err != nil {
				return fmt.Errorf("failed to delete %s %s: %w", kind, res.Name, err)
			}
		case domain.ActionNoop:
		}
	}
	return nil
}

func (k *K8sReconciler) applyResource(resource *domain.Resource) error {
	kind := extractKindName(resource.Kind)
	k.logger.Info(fmt.Sprintf("Applying %s: %s", kind, resource.Name))

	manifest, err := k.buildManifest(resource, kind)
	if err != nil {
		return fmt.Errorf("failed to build manifest for %s %s: %w", kind, resource.Name, err)
	}

	data, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "synctl-k8s-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	tmpFile.Close()

	out, err := k.runner.Run("kubectl", "apply", "-f", tmpFile.Name())
	if err != nil {
		return fmt.Errorf("kubectl apply failed: %w - output: %s", err, out.Stderr)
	}

	k.logger.Info(fmt.Sprintf("%s %s applied successfully", kind, resource.Name))
	return nil
}

func (k *K8sReconciler) buildManifest(resource *domain.Resource, kind string) (map[string]interface{}, error) {
	if kind == "raw" {
		if m, ok := resource.Spec["manifest"]; ok {
			if manifest, ok := m.(map[string]interface{}); ok {
				return manifest, nil
			}
		}
		return nil, fmt.Errorf("k8s.raw resource %s missing 'manifest' in spec", resource.Name)
	}

	apiVersion, ok := apiVersionMap[strings.ToLower(kind)]
	if !ok {
		apiVersion = "v1"
	}

	kindName := strings.ToUpper(kind[:1]) + kind[1:]

	metadata := map[string]interface{}{
		"name": resource.Name,
	}

	if !clusterScopedKinds[strings.ToLower(kind)] {
		namespace := extractNamespace(resource.Spec)
		if namespace != "" {
			metadata["namespace"] = namespace
		}
	}

	return map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kindName,
		"metadata":   metadata,
		"spec":       resource.Spec,
	}, nil
}

func (k *K8sReconciler) deleteResource(resource *domain.Resource) error {
	kind := extractKindName(resource.Kind)
	namespace := extractNamespace(resource.Spec)

	k.logger.Info(fmt.Sprintf("Deleting %s: %s", kind, resource.Name))

	args := []string{"delete", kind, resource.Name, "--ignore-not-found"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}

	out, err := k.runner.Run("kubectl", args...)
	if err != nil {
		return fmt.Errorf("kubectl delete failed: %w - output: %s", err, out.Stderr)
	}

	k.logger.Info(fmt.Sprintf("%s %s deleted successfully", kind, resource.Name))
	return nil
}

func extractKindName(kind string) string {
	parts := strings.Split(kind, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return strings.ToLower(kind)
}

func extractNamespace(spec map[string]interface{}) string {
	if ns, ok := spec["namespace"].(string); ok {
		return ns
	}
	if metadata, ok := spec["metadata"].(map[string]interface{}); ok {
		if ns, ok := metadata["namespace"].(string); ok {
			return ns
		}
	}
	return "default"
}

func isSupportedKind(kind string) bool {
	kindLower := strings.ToLower(kind)
	for _, supported := range supportedKinds {
		if supported == kindLower {
			return true
		}
	}
	return false
}
