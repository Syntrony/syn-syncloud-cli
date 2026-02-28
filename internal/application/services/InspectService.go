package services

import (
	dto "synctl/internal/domain/dto"
	interfaces "synctl/internal/domain/interfaces"
)

type InspectService struct {
	Repo interfaces.Repository
}

func (s *InspectService) Execute() (*dto.InspectDto, error) {
	exists, err := s.Repo.Exists()
	if err != nil {
		return nil, err
	}

	if !exists {
		return &dto.InspectDto{Found: false}, nil
	}

	st, err := s.Repo.Load()
	if err != nil {
		return nil, err
	}

	mode := ""

	if st.Cluster != nil {
		mode = st.Cluster.Mode
	}

	return &dto.InspectDto{
		Found:     true,
		Version:   st.Version,
		Mode:      mode,
		NodeCount: len(st.Nodes),
		ResCount:  len(st.Resources),
	}, nil

}
