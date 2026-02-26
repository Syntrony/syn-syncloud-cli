package services

import (
	dto "synctl/internal/domain/Dto"
	repository "synctl/internal/infrastructure/state"
)

type InspectService struct {
	Repo repository.Repository
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
