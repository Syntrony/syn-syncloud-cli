package services

import (
	dto "synctl/internal/domain/Dto"
	repository "synctl/internal/infrastructure/state"
)

type GetNodeService struct {
	Repo repository.Repository
}

func (s *GetNodeService) Execute() ([]dto.NodeDto, error) {
	exists, err := s.Repo.Exists()

	if err != nil {
		return nil, err
	}

	if !exists {
		return []dto.NodeDto{}, nil
	}

	st, err := s.Repo.Load()
	if err != nil {
		return nil, err
	}

	var nodes []dto.NodeDto

	for _, node := range st.Nodes {
		nodes = append(nodes, dto.NodeDto{
			Id:       node.Id,
			Hostname: node.Hostname,
			Ip:       node.Ip,
			Role:     node.Role,
		})
	}

	return nodes, nil
}
