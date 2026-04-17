package install

import (
	"synctl/internal/app"
	"synctl/internal/application/interfaces/install"
	"synctl/internal/application/services/common/banner"
	installservice "synctl/internal/application/services/install"
	dnsservice "synctl/internal/application/services/install/dns"
	dockerservice "synctl/internal/application/services/install/docker"
	k8sservice "synctl/internal/application/services/install/k8s"
	statebuilder "synctl/internal/application/services/state"

	dnsinfra "synctl/internal/infrastructure/dns"
	dockerinfra "synctl/internal/infrastructure/docker"
	k8sinfra "synctl/internal/infrastructure/k8s"
	systeminfra "synctl/internal/infrastructure/system"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install Syncloud on the cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if err := components.DetectSudo(); err != nil {
			return err
		}

		// Ensure state directory exists and is writable by the current user
		if err := components.EnsureStateDir(); err != nil {
			return err
		}

		k8sDetector := k8sinfra.NewDetector(components.Runner, components.Logger)
		k8sInspector := k8sinfra.NewInspector(components.Runner, components.Logger)
		k8sInstaller := k8sinfra.NewInstaller(components.Runner, components.Sudo, components.RuntimeType, components.Logger)

		k8sRuntime := k8sservice.NewInstallService(
			k8sDetector,
			k8sInspector,
			k8sInstaller,
		)

		dockerDetector := dockerinfra.NewDetector(components.Runner, components.Logger)
		dockerInspector := dockerinfra.NewInspector(components.Runner, components.Sudo, components.Logger)
		dockerInstaller := dockerinfra.NewInstaller(components.Runner, components.Sudo, components.RuntimeType, components.Logger)

		dockerRuntime := dockerservice.NewInstallService(
			dockerDetector,
			dockerInspector,
			dockerInstaller,
			components.Runner,
			components.Logger,
		)

		systemInspector := systeminfra.NewInspector(components.Runner)
		systemFile := systeminfra.NewFileSystem()

		dnsDetector := dnsinfra.NewDetector(components.Runner, components.Logger)
		dnsInspector := dnsinfra.NewInspector(components.Runner, components.Logger)
		dnsInstaller := dnsinfra.NewInstaller(components.Runner, components.Sudo, components.RuntimeType, components.Logger, systemInspector, systemFile)

		dnsRuntime := dnsservice.NewInstallService(
			dnsDetector,
			dnsInspector,
			dnsInstaller,
		)

		builder := statebuilder.NewStateBuilder()

		service := installservice.NewInstallService(
			components.Repo,
			[]install.RuntimeInstaller{k8sRuntime, dockerRuntime, dnsRuntime},
			systemInspector,
			builder,
			components.Logger,
		)

		err := service.Install()

		if err != nil {
			return err
		}

		banner.PrintBanner()

		return nil
	},
}
