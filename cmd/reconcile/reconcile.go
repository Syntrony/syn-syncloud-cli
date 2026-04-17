package reconcile

import (
	"synctl/internal/app"
	"synctl/internal/application/services/daemon/reconciler"
	"synctl/internal/application/services/daemon/reconciler/docker"
	dnsrecon "synctl/internal/application/services/daemon/reconciler/dns"
	"synctl/internal/application/services/daemon/reconciler/k8s"
	dnsinfra "synctl/internal/infrastructure/dns"
	sysinfra "synctl/internal/infrastructure/system"

	"synctl/internal/application/interfaces/common"

	"github.com/spf13/cobra"
)

var dryRun bool

func init() {
	Cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview reconciliation plan without applying changes")
}

var Cmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Run a one-shot reconciliation of desired vs actual state",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()
		components.DetectSudo()

		state, err := components.Repo.Load()
		if err != nil {
			return err
		}

		dockerRecon := docker.NewDockerReconciler(components.Runner, components.Logger)
		k8sRecon := k8s.NewK8sReconciler(components.Runner, components.Logger)

		runtimes := []common.RuntimeReconciler{dockerRecon, k8sRecon}

		dnsInspector := dnsinfra.NewInspector(components.Runner, components.Logger)
		dnsInstaller := dnsinfra.NewInstaller(
			components.Runner,
			components.Sudo,
			components.RuntimeType,
			components.Logger,
			components.Inspector,
			sysinfra.NewFileSystem(),
		)
		dnsRecon := dnsrecon.NewDnsReconciler(dnsInspector, dnsInstaller, components.Inspector, components.Logger)

		svc := reconciler.NewReconcileService(runtimes).WithDnsReconciler(dnsRecon)

		if dryRun {
			svc.WithDryRun(components.Logger)
			components.Logger.Info("[dry-run] Generating reconciliation plan...")
		} else {
			components.Logger.Info("Reconciling platform state...")
		}

		if err := svc.Reconcile(state); err != nil {
			return err
		}

		if dryRun {
			components.Logger.Info("[dry-run] Plan generated — no changes applied")
		} else {
			components.Logger.Info("Reconciliation completed successfully")
		}
		return nil
	},
}
