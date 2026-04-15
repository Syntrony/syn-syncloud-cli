package docker

import (
	"fmt"
	"strings"

	"synctl/internal/application/interfaces"
)

type DockerOperations struct {
	runner interfaces.CommandRunner
	logger interfaces.Logger
}

func NewDockerOperations(runner interfaces.CommandRunner, logger interfaces.Logger) *DockerOperations {
	return &DockerOperations{
		runner: runner,
		logger: logger,
	}
}

func (d *DockerOperations) ContainerExists(name string) (bool, error) {
	cmd := []string{"ps", "-a", "-f", fmt.Sprintf("name=%s", name), "--format", "{{.Names}}"}
	out, err := d.runner.Run("docker", cmd...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out.Stdout) == name, nil
}

func (d *DockerOperations) NetworkExists(name string) (bool, error) {
	cmd := []string{"network", "ls", "-f", fmt.Sprintf("name=%s", name), "--format", "{{.Name}}"}
	out, err := d.runner.Run("docker", cmd...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out.Stdout) == name, nil
}

func (d *DockerOperations) ImageExists(name string) (bool, error) {
	cmd := []string{"image", "ls", "-f", fmt.Sprintf("reference=%s", name), "--format", "{{.Repository}}"}
	out, err := d.runner.Run("docker", cmd...)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out.Stdout) == name, nil
}

func (d *DockerOperations) InspectContainer(name string) (string, error) {
	cmd := []string{"inspect", name}
	out, err := d.runner.Run("docker", cmd...)
	if err != nil {
		return "", err
	}
	return out.Stdout, nil
}

func (d *DockerOperations) RemoveContainer(name string) error {
	d.logger.Info(fmt.Sprintf("Removing container: %s", name))
	args := []string{"rm", "-f", name}
	out, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to remove container: %w - output: %s", err, out.Stderr)
	}
	d.logger.Info(fmt.Sprintf("Container %s removed", name))
	return nil
}

func (d *DockerOperations) RemoveNetwork(name string) error {
	d.logger.Info(fmt.Sprintf("Removing network: %s", name))
	args := []string{"network", "rm", name}
	out, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to remove network: %w - output: %s", err, out.Stderr)
	}
	d.logger.Info(fmt.Sprintf("Network %s removed", name))
	return nil
}

func (d *DockerOperations) RemoveImage(name string) error {
	d.logger.Info(fmt.Sprintf("Removing image: %s", name))
	args := []string{"rmi", name}
	out, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to remove image: %w - output: %s", err, out.Stderr)
	}
	d.logger.Info(fmt.Sprintf("Image %s removed", name))
	return nil
}

func (d *DockerOperations) PullImage(name string) error {
	d.logger.Info(fmt.Sprintf("Pulling image: %s", name))
	args := []string{"pull", name}
	out, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to pull image: %w - output: %s", err, out.Stderr)
	}
	d.logger.Info(fmt.Sprintf("Image %s pulled successfully", name))
	return nil
}

func (d *DockerOperations) CreateNetwork(name, driver string, labels map[string]string) error {
	d.logger.Info(fmt.Sprintf("Creating network: %s", name))
	args := []string{"network", "create"}

	if driver != "" {
		args = append(args, "-d", driver)
	}

	for k, v := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, name)

	out, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to create network: %w - output: %s", err, out.Stderr)
	}
	d.logger.Info(fmt.Sprintf("Network %s created successfully", name))
	return nil
}

func (d *DockerOperations) CreateContainer(spec map[string]interface{}, image, name string) error {
	d.logger.Info(fmt.Sprintf("Creating container: %s", name))

	args := []string{"run", "-d", "--name", name}

	if restart := GetStringValue(spec, "restart"); restart != "" {
		args = append(args, "--restart", restart)
	}

	args = d.appendPorts(args, spec)
	args = d.appendEnvVars(args, spec)
	args = d.appendVolumes(args, spec)
	args = d.appendNetwork(args, spec)

	args = append(args, image)

	out, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to create container: %w - output: %s", err, out.Stderr)
	}
	d.logger.Info(fmt.Sprintf("Container %s created successfully", name))
	return nil
}

func (d *DockerOperations) appendPorts(args []string, spec map[string]interface{}) []string {
	if ports, ok := spec["ports"].([]interface{}); ok {
		for _, p := range ports {
			if portMap, ok := p.(map[string]interface{}); ok {
				if hostPort, ok := portMap["hostPort"]; ok {
					if containerPort, ok := portMap["containerPort"]; ok {
						args = append(args, "-p", fmt.Sprintf("%v:%v", hostPort, containerPort))
					}
				}
			}
		}
	}
	return args
}

func (d *DockerOperations) appendEnvVars(args []string, spec map[string]interface{}) []string {
	if envs, ok := spec["env"].([]interface{}); ok {
		for _, e := range envs {
			if envMap, ok := e.(map[string]interface{}); ok {
				if name, ok := envMap["name"].(string); ok {
					if value, ok := envMap["value"].(string); ok {
						args = append(args, "-e", fmt.Sprintf("%s=%s", name, value))
					}
				}
			}
		}
	}
	return args
}

func (d *DockerOperations) appendVolumes(args []string, spec map[string]interface{}) []string {
	if volumes, ok := spec["volumes"].([]interface{}); ok {
		for _, v := range volumes {
			if volMap, ok := v.(map[string]interface{}); ok {
				if source, ok := volMap["source"].(string); ok {
					if target, ok := volMap["target"].(string); ok {
						args = append(args, "-v", fmt.Sprintf("%s:%s", source, target))
					}
				}
			}
		}
	}
	return args
}

func (d *DockerOperations) appendNetwork(args []string, spec map[string]interface{}) []string {
	if network, ok := spec["network"].(map[string]interface{}); ok {
		if netName, ok := network["name"].(string); ok {
			args = append(args, "--network", netName)
		}
	}
	return args
}

func GetStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func GetMapValue(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return nil
}

func GetSliceValue[T any](m map[string]interface{}, key string) []T {
	if v, ok := m[key].([]interface{}); ok {
		result := make([]T, 0, len(v))
		for _, item := range v {
			if val, ok := item.(T); ok {
				result = append(result, val)
			}
		}
		return result
	}
	return nil
}
