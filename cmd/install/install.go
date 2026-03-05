package install

import (
	"github.com/spf13/cobra"

	// infrastructure
	"synctl/internal/domain"
	executor "synctl/internal/executor"
	k8sinfra "synctl/internal/infrastructure/k8s"
	systeminfra "synctl/internal/infrastructure/system"
	loggerinfra "synctl/internal/logger"

	// services
	installservice "synctl/internal/application/services/install"
	k8sservice "synctl/internal/application/services/install/k8s"

	installiface "synctl/internal/application/interfaces/install"
)

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install Syncloud on the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {

		// =========================
		// BASE
		// =========================

		logger := loggerinfra.NewConsoleLogger()
		runner := executor.NewExecRunner(logger)
		sudo := executor.NewSudoRunner(logger)

		// =========================
		// K8S INFRA
		// =========================

		detector := k8sinfra.NewDetector(runner)
		inspector := k8sinfra.NewInspector(runner)
		installer := k8sinfra.NewInstaller(runner, sudo)

		k8sRuntime := k8sservice.NewInstallService(
			detector,
			inspector,
			installer,
		)

		// =========================
		// SYSTEM
		// =========================

		systemInspector := systeminfra.NewInspector(runner)

		builder := installservice.NewStateBuilder()

		// =========================
		// CENTRAL INSTALL
		// =========================

		service := installservice.NewInstallService(
			[]installiface.RuntimeInstaller{k8sRuntime},
			systemInspector,
			builder,
		)

		envDetector := systeminfra.NewEnvironmentDetector()

		mode := domain.Production

		if envDetector.IsContainer() {
			logger.Info("Container environment detected → using DEV mode (k3d)")
			mode = domain.Development
		}

		return service.Install(mode)
	},
}
