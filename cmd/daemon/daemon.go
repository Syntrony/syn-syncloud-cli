package daemon

import (
	"context"
	"synctl/internal/app"
	"synctl/internal/application/services/apply"
	"synctl/internal/application/services/apply/reconciler"
	"synctl/internal/application/services/apply/reconciler/docker"
	"synctl/internal/application/services/apply/reconciler/k8s"
	"time"

	"github.com/spf13/cobra"

	iApply "synctl/internal/application/interfaces/apply"
)

var Cmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run synctl control plane",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		dockerRecon := docker.NewDockerReconciler(components.Runner, components.Logger)
		k8sRecon := k8s.NewK8sReconciler(components.Runner, components.Logger)

		runtimes := []iApply.RuntimeReconciler{dockerRecon, k8sRecon}

		reconcile := reconciler.NewReconcileService(runtimes)
		controller := apply.NewController(components.Repo, reconcile, 5*time.Second, components.Logger)

		controller.Start(ctx)
		return nil
	},
}
