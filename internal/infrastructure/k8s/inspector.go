package k8s

import (
	"encoding/json"
	"strings"
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	"synctl/internal/domain/k8s"
	domainresource "synctl/internal/domain/resource"

	"github.com/google/uuid"
)

var resources = []string{
	"namespaces",
	"deployments",
	"services",
	"ingress",
	"serviceaccounts",
	"roles",
	"rolebindings",
	"clusterroles",
	"secrets",
	"endpoints",
}

type Inspector struct {
	runner interfaces.CommandRunner
	logger interfaces.Logger
}

func NewInspector(runner interfaces.CommandRunner, logger interfaces.Logger) *Inspector {
	return &Inspector{
		runner: runner,
		logger: logger,
	}
}

func (i *Inspector) Snapshot() (*domain.Snapshot, error) {
	snapshot := &domain.Snapshot{}
	i.logger.Info("Collecting Kubernetes resources...")

	for _, resource := range resources {

		i.logger.Info("Collecting " + resource + "...")

		rs, err := i.Collect([]string{}, resource, true)
		if err != nil {
			return nil, err
		}
		snapshot.Resources = append(snapshot.Resources, rs...)
	}

	return snapshot, nil

}

func (i *Inspector) Collect(
	KubectlArgs []string,
	kind string,
	namespaced bool,
) ([]domain.Resource, error) {
	resources := []domain.Resource{}

	if namespaced {
		KubectlArgs = append(KubectlArgs, "--all-namespaces")
	}

	KubectlArgs = append(KubectlArgs, "-o", "json")
	output, err := i.runner.Run("kubectl", append([]string{"get", kind}, KubectlArgs...)...)
	if err != nil {
		return nil, err
	}

	var list k8s.KubeList
	if err := json.Unmarshal([]byte(output.Stdout), &list); err != nil {
		return nil, err
	}

	for _, item := range list.Items {
		name := item.Metadata.Name
		ns := item.Metadata.Namespace
		k8sKind := strings.ToLower(item.Kind)

		// 1. Filtro de Namespaces de Sistema
		isSystemNS := ns == "kube-system" || ns == "kube-public" || ns == "kube-node-lease" || ns == "default"

		// 2. Filtro de objetos Namespace (cuando el recurso listado es un Namespace)
		isSystemNSObject := k8sKind == "namespace" && isSystemNS

		// 3. Filtro de recursos globales (ClusterRoles, etc) que son de K8s interno
		// Ignoramos los que empiezan con "system:", "admin", "edit", "view" o nombres de provisioners comunes
		isInternalRBAC := (k8sKind == "clusterrole" || k8sKind == "clusterrolebinding") && (strings.HasPrefix(name, "system:") ||
			name == "admin" || name == "edit" || name == "view" || name == "cluster-admin" ||
			strings.Contains(name, "local-path-provisioner") ||
			strings.Contains(name, "k3s-") ||
			strings.Contains(name, "traefik"))

		// 4. Filtro por nombre para recursos huérfanos en default
		isOrphanInternal := isSystemNS && (name == "kubernetes" || name == "default")

		if isSystemNS || isSystemNSObject || isInternalRBAC || isOrphanInternal {
			continue
		}

		resource := domain.Resource{
			Id:        item.Metadata.Uid,
			Name:      name,
			Kind:      "k8s." + k8sKind,
			Runtime:   "kubernetes",
			NodeId:    uuid.NewString(), // Nota: Considera usar item.Metadata.Uid para persistencia
			Source:    domain.SourceObserved,
			Spec:      item.Spec,
			Status:    item.Status,
			CreatedAt: item.Metadata.CreationTimestamp,
			UpdatedAt: item.Metadata.CreationTimestamp,
		}

		// Lógica de Ownership...
		if len(item.Metadata.OwnerReferences) > 0 {
			owner := item.Metadata.OwnerReferences[0]
			resource.Ownership = domainresource.Ownership{
				OwnerId:   owner.Uid,
				OwnerKind: owner.Kind,
				OwnerName: owner.Name,
			}
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
