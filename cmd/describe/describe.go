package describe

import (
	"fmt"

	"synctl/internal/app"
	"synctl/internal/application/outputs"
	"synctl/internal/application/services/daemon/reconciler/docker"
	"synctl/internal/application/services/daemon/reconciler/k8s"
	describeSvc "synctl/internal/application/services/describe"
	getSvc "synctl/internal/application/services/get"

	"github.com/spf13/cobra"

	"synctl/internal/application/interfaces/common"
)

var (
	name    string
	id      string
	kind    string
	runtime string
	output  string
)

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "describe a resource by name")
	Cmd.Flags().StringVarP(&id, "id", "I", "", "describe a resource by id")
	Cmd.Flags().StringVarP(&kind, "kind", "k", "", "resource kind (e.g., docker.container, k8s.deployment)")
	Cmd.Flags().StringVarP(&runtime, "runtime", "r", "", "runtime type (docker or kubernetes)")
	Cmd.Flags().StringVarP(&output, "output", "o", "", "output format (table or json)")
}

var Cmd = &cobra.Command{
	Use:   "describe -n <name> [-k <kind>] [-r <runtime>]",
	Short: "Describe Syncloud objects in wide platform scope",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if name == "" && id == "" && kind == "" {
			return fmt.Errorf("requires at least one flag: -n, -I, or -k")
		}

		dockerRecon := docker.NewDockerReconciler(components.Runner, components.Logger)
		k8sRecon := k8s.NewK8sReconciler(components.Runner, components.Logger)

		runtimes := []common.RuntimeReconciler{dockerRecon, k8sRecon}

		getService := getSvc.GetResourceService{
			Repo: components.Repo,
		}

		describeService := describeSvc.NewDescribeService(components.Repo, runtimes, getService)

		result, err := describeService.Execute(id, name, kind, runtime)
		if err != nil {
			return err
		}

		outputs.PrintDescribe(result, output)

		return nil
	},
}
