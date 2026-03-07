package install

import (
	"github.com/spf13/cobra"

	// domain
	"synctl/internal/domain"

	// infrastructure
	systemsvc "synctl/internal/application/services/system"
	"synctl/internal/domain/persistence"
	executor "synctl/internal/executor"
	k8sinfra "synctl/internal/infrastructure/k8s"
	systeminfra "synctl/internal/infrastructure/system"
	loggerinfra "synctl/internal/logger"

	// services
	installservice "synctl/internal/application/services/install"
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

		// =========================
		// BASE
		// =========================

		envDetector := systeminfra.NewEnvironmentDetector()
		logger := loggerinfra.NewConsoleLogger()
		runner := executor.NewExecRunner(logger)

		var sudo interfaces.PrivilegedRunner

		var runtime domain.RuntimeType = domain.K3s

		if envDetector.IsContainer() {
			logger.Info("Container environment detected → using DEV mode (k3d)")
			runtime = domain.K3d
		}

		if envDetector.IsRoot() {

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

		// =========================
		// K8S INFRA
		// =========================

		detector := k8sinfra.NewDetector(runner)
		inspector := k8sinfra.NewInspector(runner)
		installer := k8sinfra.NewInstaller(runner, sudo, runtime)

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
			repo,
			[]installiface.RuntimeInstaller{k8sRuntime},
			systemInspector,
			builder,
		)

		return service.Install()
	},
}
