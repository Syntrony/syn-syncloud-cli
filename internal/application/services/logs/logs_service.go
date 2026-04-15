package logs

import (
	"fmt"
	"strings"
	"synctl/internal/application/interfaces"
)

type LogsService struct {
	runner interfaces.CommandRunner
}

func NewLogsService(runner interfaces.CommandRunner) *LogsService {
	return &LogsService{runner: runner}
}

type LogsOptions struct {
	Name      string
	Runtime   string
	Namespace string
	Follow    bool
	Tail      int
}

func (s *LogsService) Fetch(opts LogsOptions) (string, error) {
	runtime := strings.ToLower(opts.Runtime)

	switch runtime {
	case "docker":
		return s.dockerLogs(opts)
	case "kubernetes", "k8s":
		return s.k8sLogs(opts)
	default:
		return "", fmt.Errorf("unsupported runtime '%s': use 'docker' or 'kubernetes'", opts.Runtime)
	}
}

func (s *LogsService) dockerLogs(opts LogsOptions) (string, error) {
	args := []string{"logs"}
	if opts.Follow {
		args = append(args, "-f")
	}
	if opts.Tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", opts.Tail))
	}
	args = append(args, opts.Name)

	out, err := s.runner.Run("docker", args...)
	if err != nil {
		return "", fmt.Errorf("docker logs failed: %w", err)
	}
	return out.Stdout, nil
}

func (s *LogsService) k8sLogs(opts LogsOptions) (string, error) {
	args := []string{"logs"}
	if opts.Follow {
		args = append(args, "-f")
	}
	if opts.Tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", opts.Tail))
	}
	if opts.Namespace != "" {
		args = append(args, "-n", opts.Namespace)
	}
	args = append(args, opts.Name)

	out, err := s.runner.Run("kubectl", args...)
	if err != nil {
		return "", fmt.Errorf("kubectl logs failed: %w", err)
	}
	return out.Stdout, nil
}
