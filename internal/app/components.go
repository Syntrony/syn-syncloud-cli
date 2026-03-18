package app

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/domain"
	"synctl/internal/domain/persistence"
	executor "synctl/internal/executor"
	"synctl/internal/infrastructure/system"
	"synctl/internal/logger"
)

const DefaultStatePath = "helpers/syncloud-state.json"

type Components struct {
	Repo        interfaces.Repository
	Logger      interfaces.Logger
	Runner      interfaces.CommandRunner
	Sudo        interfaces.PrivilegedRunner
	EnvDetector *system.EnvironmentDetector
	RuntimeType domain.RuntimeType
}

func NewComponents() *Components {
	return &Components{}
}

func (c *Components) Init() {
	c.Logger = logger.NewConsoleLogger()
	c.Runner = executor.NewExecRunner(c.Logger)
	c.EnvDetector = system.NewEnvironmentDetector()
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
	return c.WithRepo(DefaultStatePath)
}

func (c *Components) WithSudo(sudo interfaces.PrivilegedRunner) *Components {
	c.Sudo = sudo
	return c
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
