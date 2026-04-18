package services

import (
	interfaces "synctl/internal/application/interfaces"
	dto "synctl/internal/domain/dto"
)

type GetDnsService struct {
	Repo interfaces.Repository
}

func (s *GetDnsService) Execute() ([]dto.DnsRecordDto, error) {
	exists, err := s.Repo.Exists()
	if err != nil {
		return nil, err
	}

	if !exists {
		return []dto.DnsRecordDto{}, nil
	}

	st, err := s.Repo.Load()
	if err != nil {
		return nil, err
	}

	if st.Dns == nil || st.Dns.Records == nil {
		return []dto.DnsRecordDto{}, nil
	}

	var records []dto.DnsRecordDto
	for _, r := range st.Dns.Records {
		records = append(records, dto.DnsRecordDto{
			Name:   r.Name,
			Server: r.Server,
		})
	}

	return records, nil
}
