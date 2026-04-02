package dns

import (
	"fmt"
	"strings"
	"synctl/internal/application/interfaces"
	system "synctl/internal/application/interfaces/system"
	"synctl/internal/domain"
	resource "synctl/internal/domain/Resource"

	"github.com/brianvoe/gofakeit/v6"
)

type Installer struct {
	runner       interfaces.CommandRunner
	sudo         interfaces.PrivilegedRunner
	runtime      domain.RuntimeType
	logger       interfaces.Logger
	sysinspector system.Inspector
	sysfile      system.FileSystem
}

func NewInstaller(runner interfaces.CommandRunner,
	sudo interfaces.PrivilegedRunner, runtime domain.RuntimeType,
	logger interfaces.Logger, sysinspector system.Inspector, sysfile system.FileSystem) *Installer {
	return &Installer{
		runner:       runner,
		sudo:         sudo,
		runtime:      runtime,
		logger:       logger,
		sysinspector: sysinspector,
		sysfile:      sysfile,
	}
}

func (i *Installer) InstallDns() error {

	if i.runtime == domain.Container {
		i.logger.Info("Skipping dns install inside container runtime")
		return nil
	}

	i.logger.Info("Installing dnsmasq...")

	_, err := i.sudo.Run("apt-get", "install", "-y", "dnsmasq")

	if err != nil {
		return fmt.Errorf("dnsmasq install failed: %w", err)
	}

	if err = i.DisableSystemd(); err != nil {
		return err
	}

	_, err = i.sudo.Run("systemctl", "enable", "dnsmasq")

	if err != nil {
		return err
	}

	return nil
}

// Revisar pipeline del repositorio de Networking
func (i *Installer) ConfigureDns() error {
	i.logger.Info("Configuring dnsmasq file...")

	randomWord := gofakeit.Word()

	host, err := i.sysinspector.Inspect()

	if err != nil {
		return err
	}

	record := resource.Record{
		Name:   randomWord + ".syncloud.local",
		Server: host[0].Ip,
	}

	content := i.SetupDnsFile([]resource.Record{record}, host[0].Ip)

	// Definir rutas
	finalPath := "/etc/dnsmasq.d/syncloud.conf"
	tmpPath := "/tmp/syncloud.conf.tmp"
	backupPath := fmt.Sprintf("/etc/dnsmasq.d/syncloud.conf.bak")

	// Escribir archivo temporal
	if err := i.sysfile.WriteFile(tmpPath, content); err != nil {
		return fmt.Errorf("failed to write tmp config: %w", err)
	}

	// Validacion de archivo temporal
	_, err = i.runner.Run("dnsmasq", "--test", "--conf-file="+tmpPath)

	if err != nil {
		return fmt.Errorf("invalid dnsmasq configuration: %w", err)
	}

	if i.runtime == domain.Container {
		// Backup del archivo
		i.runner.Run("cp", finalPath, backupPath)

		if _, err := i.runner.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("Failed to deploy config: %w", err)
		}

	} else {
		i.sudo.Run("cp", finalPath, backupPath)

		if _, err := i.sudo.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("Failed to deploy config: %w", err)
		}

		_, err = i.sudo.Run("systemctl", "reload", "dnsmasq")

		if err != nil {
			i.logger.Error("Reload failed, performing rollback...")
			i.sudo.Run("cp", backupPath, finalPath)
			i.sudo.Run("systemctl", "reload", "dnsmasq")
			return fmt.Errorf("dnsmasq reload failed, rolled back: %w", err)
		}
	}

	i.logger.Info("DNS config applied successfully with record: " + record.Name)
	return nil
}

// ApplyRecords writes the dnsmasq config with the given records and reloads the service.
func (i *Installer) ApplyRecords(records []resource.Record, generalIp string) error {
	content := i.SetupDnsFile(records, generalIp)

	finalPath := "/etc/dnsmasq.d/syncloud.conf"
	tmpPath := "/tmp/syncloud.conf.tmp"
	backupPath := "/etc/dnsmasq.d/syncloud.conf.bak"

	if err := i.sysfile.WriteFile(tmpPath, content); err != nil {
		return fmt.Errorf("failed to write tmp config: %w", err)
	}

	if _, err := i.runner.Run("dnsmasq", "--test", "--conf-file="+tmpPath); err != nil {
		return fmt.Errorf("invalid dnsmasq configuration: %w", err)
	}

	if i.runtime == domain.Container {
		i.runner.Run("cp", finalPath, backupPath)
		if _, err := i.runner.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("failed to deploy config: %w", err)
		}
	} else {
		i.sudo.Run("cp", finalPath, backupPath)
		if _, err := i.sudo.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("failed to deploy config: %w", err)
		}
		if _, err := i.sudo.Run("systemctl", "reload", "dnsmasq"); err != nil {
			i.logger.Error("Reload failed, performing rollback...")
			i.sudo.Run("cp", backupPath, finalPath)
			i.sudo.Run("systemctl", "reload", "dnsmasq")
			return fmt.Errorf("dnsmasq reload failed, rolled back: %w", err)
		}
	}

	i.logger.Info(fmt.Sprintf("DNS config applied with %d record(s)", len(records)))
	return nil
}

func (i *Installer) DisableSystemd() error {
	i.logger.Info("Disabling systemd...")

	_, err := i.sudo.Run("systemctl", "disable", "systemd-resolved")

	if err != nil {
		return fmt.Errorf("systemd-resolved disable failed: %w", err)
	}

	return nil
}

func (i *Installer) SetupDnsFile(records []resource.Record, generalIp string) string {
	var b strings.Builder

	b.WriteString("# ==========================================\n")
	b.WriteString("# SYNCLOUD DNS CONFIGURATION\n")
	b.WriteString("# ==========================================\n")

	for _, r := range records {
		b.WriteString(fmt.Sprintf("address=/%s/%s\n",
			r.Name, r.Server))
	}

	b.WriteString("\nserver=8.8.8.8\n")

	b.WriteString("server=8.8.4.4\n")

	b.WriteString("listen-address=127.0.0.1," + generalIp + "\n")

	b.WriteString("bind-interfaces\n")

	return b.String()
}
