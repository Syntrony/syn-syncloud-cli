package dns

import (
	"synctl/internal/application/interfaces/install"
	install_dns "synctl/internal/application/interfaces/install/dns"
	"synctl/internal/domain"
)

type InstallService struct {
	detector  install.Detector
	inspector install.Inspector
	installer install_dns.DnsInstaller
}

func NewInstallService(detector install.Detector, inspector install.Inspector, installer install_dns.DnsInstaller) *InstallService {
	return &InstallService{
		detector:  detector,
		inspector: inspector,
		installer: installer,
	}
}

func (s *InstallService) Install() (*domain.Snapshot, error) {
	status, err := s.detector.Detect()

	if err != nil {
		return nil, err
	}

	if !status.DnsInstalled {
		if err := s.installer.InstallDns(); err != nil {
			return nil, err
		}
	}

	if !status.DnsConfigured {
		if err := s.installer.ConfigureDns(); err != nil {
			return nil, err
		}
	}

	return s.inspector.Snapshot()
}
