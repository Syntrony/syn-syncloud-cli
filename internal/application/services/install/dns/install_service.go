package dns

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"synctl/internal/application/interfaces/install"
	install_dns "synctl/internal/application/interfaces/install/dns"
	"synctl/internal/domain"
)

const syncloudDefaultConfDir = "/etc/dnsmasq.d"

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
		// Detect dnsmasq conf directories that exist outside syncloud's default.
		externalDir := detectExternalConfDir()
		if externalDir != "" {
			chosen := promptConfDirChoice(externalDir)
			s.installer.WithConfDir(chosen)
		}

		if err := s.installer.ConfigureDns(); err != nil {
			return nil, err
		}
	}

	return s.inspector.Snapshot()
}

// detectExternalConfDir returns the first dnsmasq conf-dir found in /etc/dnsmasq.conf
// that is NOT syncloud's default directory. Returns "" when none is detected.
func detectExternalConfDir() string {
	f, err := os.Open("/etc/dnsmasq.conf")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "conf-dir=") {
			continue
		}
		raw := strings.TrimPrefix(line, "conf-dir=")
		// conf-dir= accepts optional suffix filters: conf-dir=/path,.conf
		dir := strings.TrimSpace(strings.SplitN(raw, ",", 2)[0])
		dir = filepath.Clean(dir)
		if dir != "" && dir != syncloudDefaultConfDir {
			return dir
		}
	}
	return ""
}

// promptConfDirChoice prints an interactive prompt when an external dnsmasq
// config directory is detected and returns the directory syncloud should use.
func promptConfDirChoice(externalDir string) string {
	const reset = "\033[0m"
	const cyan = "\033[1;36m"
	const bold = "\033[1m"

	fmt.Println()
	fmt.Println(cyan + "  ┌─────────────────────────────────────────────────────────────────┐" + reset)
	fmt.Println(cyan + "  │              ⚙  CONFIGURACIÓN DNSMASQ DETECTADA                 │" + reset)
	fmt.Println(cyan + "  ├─────────────────────────────────────────────────────────────────┤" + reset)
	fmt.Printf(cyan+"  │  Se encontraron configuraciones en: "+bold+"%-29s"+reset+cyan+"│\n"+reset, externalDir)
	fmt.Println(cyan + "  │                                                                 │" + reset)
	fmt.Println(cyan + "  │  Selecciona una opción:                                         │" + reset)
	fmt.Println(cyan + "  │                                                                 │" + reset)
	fmt.Printf(cyan+"  │  "+bold+"[1]"+reset+cyan+" Usar directorio syncloud (%-36s│\n"+reset, syncloudDefaultConfDir+")")
	fmt.Println(cyan + "  │      Actualiza dnsmasq para incluir el directorio syncloud.     │" + reset)
	fmt.Println(cyan + "  │                                                                 │" + reset)
	fmt.Printf(cyan+"  │  "+bold+"[2]"+reset+cyan+" Usar tu directorio existente (%-33s│\n"+reset, externalDir+")")
	fmt.Println(cyan + "  │      Syncloud respetará y escribirá en ese directorio.          │" + reset)
	fmt.Println(cyan + "  └─────────────────────────────────────────────────────────────────┘" + reset)
	fmt.Println()
	fmt.Print("  Ingresa tu opción [1/2] (default: 1): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return syncloudDefaultConfDir
	}

	switch strings.TrimSpace(input) {
	case "2":
		fmt.Println()
		fmt.Printf("  [OK] Syncloud usará el directorio: %s\n\n", externalDir)
		return externalDir
	default:
		fmt.Println()
		fmt.Printf("  [OK] Syncloud usará el directorio: %s\n\n", syncloudDefaultConfDir)
		return syncloudDefaultConfDir
	}
}
