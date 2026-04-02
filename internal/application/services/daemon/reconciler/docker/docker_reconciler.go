package docker

import (
	"strings"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
)

type DockerReconciler struct {
	ops    *DockerOperations
	logger interfaces.Logger
}

func NewDockerReconciler(runner interfaces.CommandRunner, logger interfaces.Logger) *DockerReconciler {
	return &DockerReconciler{
		ops:    NewDockerOperations(runner, logger),
		logger: logger,
	}
}

func (d *DockerReconciler) Runtime() string {
	return "docker"
}

func (d *DockerReconciler) Observe(resource *domain.Resource) (*domain.RuntimeResourceState, error) {
	out, err := d.ops.InspectContainer(resource.Name)
	if err != nil || out == "" {
		return &domain.RuntimeResourceState{Exists: false}, nil
	}

	return &domain.RuntimeResourceState{
		Exists: true,
		Spec:   map[string]interface{}{"raw": out},
	}, nil
}

func (d *DockerReconciler) Reconcile(ctx domain.ReconcileContext) error {
	for _, action := range ctx.Actions {
		res := action.Resource
		kind := extractKind(res.Kind)

		switch kind {
		case "container":
			if err := d.reconcileContainer(res, action.Type); err != nil {
				return err
			}
		case "network":
			if err := d.reconcileNetwork(res, action.Type); err != nil {
				return err
			}
		case "image":
			if err := d.reconcileImage(res, action.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DockerReconciler) reconcileContainer(res *domain.Resource, actionType domain.ActionType) error {
	switch actionType {
	case domain.ActionCreate:
		exists, err := d.ops.ContainerExists(res.Name)
		if err != nil {
			return err
		}
		if exists {
			d.logger.Info("Container " + res.Name + " already exists, skipping")
			return nil
		}
		image := GetStringValue(res.Spec, "image")
		if image == "" {
			return nil
		}
		return d.ops.CreateContainer(res.Spec, image, res.Name)
	case domain.ActionDelete:
		return d.ops.RemoveContainer(res.Name)
	case domain.ActionNoop:
		return nil
	}
	return nil
}

func (d *DockerReconciler) reconcileNetwork(res *domain.Resource, actionType domain.ActionType) error {
	switch actionType {
	case domain.ActionCreate:
		exists, err := d.ops.NetworkExists(res.Name)
		if err != nil {
			return err
		}
		if exists {
			d.logger.Info("Network " + res.Name + " already exists, skipping")
			return nil
		}
		driver := GetStringValue(res.Spec, "driver")
		labels := extractLabels(res.Spec)
		return d.ops.CreateNetwork(res.Name, driver, labels)
	case domain.ActionDelete:
		return d.ops.RemoveNetwork(res.Name)
	case domain.ActionNoop:
		return nil
	}
	return nil
}

func (d *DockerReconciler) reconcileImage(res *domain.Resource, actionType domain.ActionType) error {
	switch actionType {
	case domain.ActionCreate:
		imageName := GetStringValue(res.Spec, "name")
		if imageName == "" {
			imageName = res.Name
		}
		exists, err := d.ops.ImageExists(imageName)
		if err != nil {
			return err
		}
		if exists {
			d.logger.Info("Image " + imageName + " already exists, skipping")
			return nil
		}
		return d.ops.PullImage(imageName)
	case domain.ActionDelete:
		imageName := GetStringValue(res.Spec, "name")
		if imageName == "" {
			imageName = res.Name
		}
		return d.ops.RemoveImage(imageName)
	case domain.ActionNoop:
		return nil
	}
	return nil
}

func extractKind(kind string) string {
	parts := strings.Split(kind, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return strings.ToLower(kind)
}

func extractLabels(spec map[string]interface{}) map[string]string {
	labels := make(map[string]string)
	if labelMap, ok := spec["labels"].(map[string]interface{}); ok {
		for k, v := range labelMap {
			if strVal, ok := v.(string); ok {
				labels[k] = strVal
			}
		}
	}
	return labels
}
