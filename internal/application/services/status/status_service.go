package status

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain/dto"
)

type StatusService struct {
	repo   interfaces.Repository
	runner interfaces.CommandRunner
}

func NewStatusService(repo interfaces.Repository, runner interfaces.CommandRunner) *StatusService {
	return &StatusService{repo: repo, runner: runner}
}

func (s *StatusService) Execute() (*dto.StatusDto, error) {
	result := &dto.StatusDto{
		Docker:     s.checkDocker(),
		Kubernetes: s.checkKubernetes(),
	}

	if exists, _ := s.repo.Exists(); exists {
		state, err := s.repo.Load()
		if err != nil {
			return result, nil
		}
		result.StateFound = true
		result.Version = state.Version
		result.NodeCount = len(state.Nodes)
		result.ResourceCount = len(state.Resources)
		if state.Cluster != nil {
			result.ClusterName = state.Cluster.Name
			result.ClusterMode = state.Cluster.Mode
		}
	}

	return result, nil
}

func (s *StatusService) checkDocker() string {
	out, err := s.runner.Run("docker", "info", "--format", "{{.ServerVersion}}")
	if err != nil || out.Stdout == "" {
		return "unavailable"
	}
	v := trim(out.Stdout)
	if v == "" {
		return "unavailable"
	}
	return "running (v" + v + ")"
}

func (s *StatusService) checkKubernetes() string {
	out, err := s.runner.Run("kubectl", "version", "--client", "--short")
	if err != nil {
		// Try without --short (newer kubectl)
		out, err = s.runner.Run("kubectl", "version", "--client")
		if err != nil {
			return "unavailable"
		}
	}
	v := trim(out.Stdout)
	if v == "" {
		return "unavailable"
	}
	return "available"
}

func trim(s string) string {
	result := ""
	for _, c := range s {
		if c == '\n' || c == '\r' {
			break
		}
		result += string(c)
	}
	return result
}
