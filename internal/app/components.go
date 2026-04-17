package app

import (
	"fmt"
	"path/filepath"

	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	"synctl/internal/domain/persistence"
	executor "synctl/internal/executor"
	"synctl/internal/infrastructure/system"
	"synctl/internal/logger"
)

const (
	ContainerStatePath = "/helpers/syncloud-state.json"
	VPSStatePath       = "/var/lib/synctl/syncloud-state.json"
)

type Components struct {
	Repo        interfaces.Repository
	Logger      interfaces.Logger
	Runner      interfaces.CommandRunner
	Sudo        interfaces.PrivilegedRunner
	EnvDetector *system.EnvironmentDetector
	Inspector   *system.Inspector
	RuntimeType domain.RuntimeType
}

func NewComponents() *Components {
	return &Components{}
}

func (c *Components) Init() {
	c.Logger = logger.NewConsoleLogger()
	c.Runner = executor.NewExecRunner(c.Logger)
	c.EnvDetector = system.NewEnvironmentDetector()
	c.Inspector = system.NewInspector(c.Runner)
	c.detectEnvironment()
}

func (c *Components) detectEnvironment() {
	if c.EnvDetector.IsContainer() {
		c.Logger.Info("Container environment detected -> using DEV mode (k3d)")
		c.RuntimeType = domain.Container
	} else {
		c.RuntimeType = domain.VPS
	}
}

func (c *Components) WithRepo(path string) *Components {
	c.Repo = &persistence.StateRepository{Path: path}
	return c
}

func (c *Components) WithDefaultRepo() *Components {
	if c.EnvDetector.IsContainer() {
		return c.WithRepo(ContainerStatePath)
	}
	return c.WithRepo(VPSStatePath)
}

func (c *Components) WithSudo(sudo interfaces.PrivilegedRunner) *Components {
	c.Sudo = sudo
	return c
}

// EnsureStateDir creates the VPS state directory via sudo and transfers ownership
// to the current user, so the process can write the state file without root.
func (c *Components) EnsureStateDir() error {
	if c.RuntimeType != domain.VPS {
		return nil
	}
	dir := filepath.Dir(VPSStatePath)
	if _, err := c.Sudo.Run("mkdir", "-p", dir); err != nil {
		return fmt.Errorf("failed to create state directory %s: %w", dir, err)
	}
	user, err := c.Runner.Whoami()
	if err != nil {
		return fmt.Errorf("failed to determine current user: %w", err)
	}
	if _, err := c.Sudo.Run("chown", user+":"+user, dir); err != nil {
		// Non-fatal: directory may already be owned correctly
		c.Logger.Info("chown state dir skipped: " + err.Error())
	}
	return nil
}

func (c *Components) DetectSudo() error {
	if c.EnvDetector.IsRoot() {
		c.Logger.Info("Running as root -> sudo not required")
		c.Sudo = c.Runner
		return nil
	}

	c.Logger.Info("Running as user -> sudo required")
	c.Sudo = executor.NewSudoRunner(c.Runner)
	return nil
}
