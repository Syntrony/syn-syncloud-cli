package dns

import (
	"fmt"
	"path/filepath"
	"strings"
	"synctl/internal/application/interfaces"
	system "synctl/internal/application/interfaces/system"
	"synctl/internal/domain"
	resource "synctl/internal/domain/resource"

	"github.com/brianvoe/gofakeit/v6"
)

const defaultConfDir = "/etc/dnsmasq.d"

type Installer struct {
	runner       interfaces.CommandRunner
	sudo         interfaces.PrivilegedRunner
	runtime      domain.RuntimeType
	logger       interfaces.Logger
	sysinspector system.Inspector
	sysfile      system.FileSystem
	confDir      string
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
		confDir:      defaultConfDir,
	}
}

// WithConfDir overrides the directory where syncloud writes its dnsmasq config file.
func (i *Installer) WithConfDir(dir string) {
	i.confDir = dir
}

// syncloudConfPath returns the absolute path of the syncloud dnsmasq config file.
func (i *Installer) syncloudConfPath() string {
	return filepath.Join(i.confDir, "syncloud.conf")
}

// syncloudBackupPath returns the absolute path of the syncloud dnsmasq config backup.
func (i *Installer) syncloudBackupPath() string {
	return filepath.Join(i.confDir, "syncloud.conf.bak")
}

// ensureConfDir creates the conf directory if it does not exist, and ensures
// /etc/dnsmasq.conf has a conf-dir= directive pointing to it so dnsmasq reads the files.
func (i *Installer) ensureConfDir() error {
	if i.runtime == domain.Container {
		if _, err := i.runner.Run("mkdir", "-p", i.confDir); err != nil {
			return fmt.Errorf("failed to create conf dir %s: %w", i.confDir, err)
		}
	} else {
		if _, err := i.sudo.Run("mkdir", "-p", i.confDir); err != nil {
			return fmt.Errorf("failed to create conf dir %s: %w", i.confDir, err)
		}
		// Ensure /etc/dnsmasq.conf includes the conf-dir directive so dnsmasq reads syncloud.conf
		if err := i.ensureConfDirDirective(); err != nil {
			i.logger.Info("Warning: could not update /etc/dnsmasq.conf with conf-dir: " + err.Error())
		}
	}
	return nil
}

// ensureConfDirDirective appends conf-dir=<confDir> to /etc/dnsmasq.conf if not already present.
func (i *Installer) ensureConfDirDirective() error {
	const mainConf = "/etc/dnsmasq.conf"
	directive := "conf-dir=" + i.confDir

	// Check if the directive already exists (active, not commented out)
	out, _ := i.runner.Run("grep", "-E", "^conf-dir="+i.confDir, mainConf)
	if out != nil && strings.TrimSpace(out.Stdout) != "" {
		return nil // already present
	}

	i.logger.Info("Adding conf-dir directive to " + mainConf + "...")
	line := "\n# Added by syncloud\n" + directive + "\n"
	_, err := i.sudo.Run("sh", "-c", fmt.Sprintf("echo '%s' >> %s", line, mainConf))
	return err
}

func (i *Installer) InstallDns() error {

	if i.runtime == domain.Container {
		i.logger.Info("Skipping dns install inside container runtime")
		return nil
	}

	i.logger.Info("Installing dnsmasq...")

	// DEBIAN_FRONTEND=noninteractive prevents apt-get from hanging on prompts
	_, err := i.sudo.Run("sh", "-c", "DEBIAN_FRONTEND=noninteractive apt-get install -y dnsmasq")

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

	// Start the service so it is running before ConfigureDns tries to reload it
	if _, err = i.sudo.Run("systemctl", "start", "dnsmasq"); err != nil {
		i.logger.Info("dnsmasq start failed (may already be running): " + err.Error())
	}

	return nil
}

// Revisar pipeline del repositorio de Networking
func (i *Installer) ConfigureDns() error {
	i.logger.Info("Configuring dnsmasq file...")

	if err := i.ensureConfDir(); err != nil {
		return err
	}

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
	finalPath := i.syncloudConfPath()
	tmpPath := "/tmp/syncloud.conf.tmp"
	backupPath := i.syncloudBackupPath()

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
		// Backup only if the file already exists (first run has no prior config)
		i.runner.Run("sh", "-c", fmt.Sprintf("test -f %s && cp %s %s", finalPath, finalPath, backupPath))

		if _, err := i.runner.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("Failed to deploy config: %w", err)
		}

	} else {
		// Backup only if the file already exists (first run has no prior config)
		i.sudo.Run("sh", "-c", fmt.Sprintf("test -f %s && cp %s %s", finalPath, finalPath, backupPath))

		if _, err := i.sudo.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("Failed to deploy config: %w", err)
		}

		_, err = i.sudo.Run("systemctl", "restart", "dnsmasq")

		if err != nil {
			i.logger.Error("Restart failed, performing rollback...")
			i.sudo.Run("cp", backupPath, finalPath)
			i.sudo.Run("systemctl", "restart", "dnsmasq")
			return fmt.Errorf("dnsmasq restart failed, rolled back: %w", err)
		}
	}

	i.logger.Info("DNS config applied successfully with record: " + record.Name)
	return nil
}

// ApplyRecords writes the dnsmasq config with the given records and reloads the service.
func (i *Installer) ApplyRecords(records []resource.Record, generalIp string) error {
	if err := i.ensureConfDir(); err != nil {
		return err
	}

	content := i.SetupDnsFile(records, generalIp)

	finalPath := i.syncloudConfPath()
	tmpPath := "/tmp/syncloud.conf.tmp"
	backupPath := i.syncloudBackupPath()

	if err := i.sysfile.WriteFile(tmpPath, content); err != nil {
		return fmt.Errorf("failed to write tmp config: %w", err)
	}

	if _, err := i.runner.Run("dnsmasq", "--test", "--conf-file="+tmpPath); err != nil {
		return fmt.Errorf("invalid dnsmasq configuration: %w", err)
	}

	if i.runtime == domain.Container {
		i.runner.Run("sh", "-c", fmt.Sprintf("test -f %s && cp %s %s", finalPath, finalPath, backupPath))
		if _, err := i.runner.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("failed to deploy config: %w", err)
		}
	} else {
		i.sudo.Run("sh", "-c", fmt.Sprintf("test -f %s && cp %s %s", finalPath, finalPath, backupPath))
		if _, err := i.sudo.Run("mv", tmpPath, finalPath); err != nil {
			return fmt.Errorf("failed to deploy config: %w", err)
		}
		if _, err := i.sudo.Run("systemctl", "restart", "dnsmasq"); err != nil {
			i.logger.Error("Restart failed, performing rollback...")
			i.sudo.Run("cp", backupPath, finalPath)
			i.sudo.Run("systemctl", "restart", "dnsmasq")
			return fmt.Errorf("dnsmasq restart failed, rolled back: %w", err)
		}
	}

	i.logger.Info(fmt.Sprintf("DNS config applied with %d record(s)", len(records)))
	return nil
}

func (i *Installer) DisableSystemd() error {
	i.logger.Info("Disabling systemd-resolved...")

	if _, err := i.sudo.Run("systemctl", "disable", "--now", "systemd-resolved"); err != nil {
		// May already be disabled or not present — log and continue
		i.logger.Info("systemd-resolved disable skipped (may not be active): " + err.Error())
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
