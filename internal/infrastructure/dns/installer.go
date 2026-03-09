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

	content := i.SetupDnsFile([]resource.Record{record})

	return i.sysfile.WriteFile(
		"/etc/dnsmasq.d/syncloud.conf",
		content,
	)
}

func (i *Installer) DisableSystemd() error {
	i.logger.Info("Disabling systemd...")

	_, err := i.sudo.Run("systemctl", "disable", "systemd-resolved")

	if err != nil {
		return fmt.Errorf("systemd-resolved disable failed: %w", err)
	}

	return nil
}

func (i *Installer) SetupDnsFile(records []resource.Record) string {
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

	return b.String()
}
