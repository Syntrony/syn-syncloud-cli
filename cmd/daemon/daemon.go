package daemon

import (
	"context"
	"synctl/internal/app"
	"synctl/internal/application/services/daemon"
	"synctl/internal/application/services/daemon/reconciler"
	"synctl/internal/application/services/daemon/reconciler/docker"
	"synctl/internal/application/services/daemon/reconciler/k8s"
	"time"

	"github.com/spf13/cobra"

	"synctl/internal/application/interfaces/common"
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

		runtimes := []common.RuntimeReconciler{dockerRecon, k8sRecon}

		reconcile := reconciler.NewReconcileService(runtimes)
		controller := daemon.NewController(components.Repo, reconcile, 5*time.Second, components.Logger)

		controller.Start(ctx)
		return nil
	},
}
