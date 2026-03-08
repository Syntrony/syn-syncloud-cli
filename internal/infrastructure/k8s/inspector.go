package k8s

import (
	"encoding/json"
	"strings"
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	domainresource "synctl/internal/domain/Resource"
	"synctl/internal/domain/k8s"
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

		resource := domain.Resource{
			Id:      item.Metadata.Uid,
			Name:    item.Metadata.Name,
			Kind:    "k8s." + strings.ToLower(item.Kind),
			Runtime: "kubernetes",
			NodeId:  "",
			Spec: domainresource.Spec{
				Raw: item.Spec,
			},
			Status: domainresource.Status{
				Raw: item.Status,
			},
			CreatedAt: item.Metadata.CreationTimestamp,
			UpdatedAt: item.Metadata.CreationTimestamp,
		}

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
