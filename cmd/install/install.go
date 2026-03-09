package install

import (
	"github.com/spf13/cobra"

	// domain
	"synctl/internal/domain"

	// infrastructure
	systemsvc "synctl/internal/application/services/system"
	"synctl/internal/domain/persistence"
	executor "synctl/internal/executor"
	dnsinfra "synctl/internal/infrastructure/dns"
	dockerinfra "synctl/internal/infrastructure/docker"
	k8sinfra "synctl/internal/infrastructure/k8s"
	systeminfra "synctl/internal/infrastructure/system"
	loggerinfra "synctl/internal/logger"

	// services
	installservice "synctl/internal/application/services/install"
	dnsservice "synctl/internal/application/services/install/dns"
	dockerservice "synctl/internal/application/services/install/docker"
	k8sservice "synctl/internal/application/services/install/k8s"

	interfaces "synctl/internal/application/interfaces"
	installiface "synctl/internal/application/interfaces/install"
)

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install Syncloud on the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := &persistence.StateRepository{
			Path: "helpers/syncloud-state.json",
		}

		envDetector := systeminfra.NewEnvironmentDetector()
		logger := loggerinfra.NewConsoleLogger()
		runner := executor.NewExecRunner(logger)

		var sudo interfaces.PrivilegedRunner

		var runtime domain.RuntimeType = domain.VPS

		if envDetector.IsContainer() {
			logger.Info("Container environment detected -> using DEV mode (k3d)")
			runtime = domain.Container
			sudo = runner
		} else if envDetector.IsRoot() {
			logger.Info("Running as root -> sudo not required")
			sudo = runner
		} else {
			logger.Info("Running as user -> sudo required")
			sudoSession := &systemsvc.SudoSession{}

			if err := sudoSession.Ensure(runner); err != nil {
				return err
			}

			sudo = executor.NewSudoRunner(runner)
		}

		k8sDetector := k8sinfra.NewDetector(runner, logger)
		k8sInspector := k8sinfra.NewInspector(runner, logger)
		k8sInstaller := k8sinfra.NewInstaller(runner, sudo, runtime, logger)

		k8sRuntime := k8sservice.NewInstallService(
			k8sDetector,
			k8sInspector,
			k8sInstaller,
		)

		dockerDetector := dockerinfra.NewDetector(runner, logger)
		dockerInspector := dockerinfra.NewInspector(runner, logger)
		dockerInstaller := dockerinfra.NewInstaller(runner, sudo, runtime, logger)

		dockerRuntime := dockerservice.NewInstallService(
			dockerDetector,
			dockerInspector,
			dockerInstaller,
		)

		systemInspector := systeminfra.NewInspector(runner)

		systemFile := systeminfra.NewFileSystem()

		dnsDetector := dnsinfra.NewDetector(runner, logger)
		dnsInspector := dnsinfra.NewInspector(runner, logger)
		dnsInstaller := dnsinfra.NewInstaller(runner, sudo, runtime, logger, systemInspector, systemFile)

		dnsRuntime := dnsservice.NewInstallService(
			dnsDetector,
			dnsInspector,
			dnsInstaller,
		)

		builder := installservice.NewStateBuilder()

		service := installservice.NewInstallService(
			repo,
			[]installiface.RuntimeInstaller{k8sRuntime, dockerRuntime, dnsRuntime},
			systemInspector,
			builder,
			logger,
		)

		return service.Install()
	},
}
