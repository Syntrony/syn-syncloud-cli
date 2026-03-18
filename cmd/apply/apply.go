package apply

import (
	"os"

	"synctl/internal/app"
	iApply "synctl/internal/application/interfaces/apply"
	applyService "synctl/internal/application/services/apply"
	"synctl/internal/application/services/apply/parser"
	"synctl/internal/application/services/apply/reconciler/docker"
	"synctl/internal/application/services/apply/reconciler/k8s"
	validatorsvc "synctl/internal/application/services/apply/validator"
	"synctl/internal/application/services/state"

	"github.com/spf13/cobra"
)

var filePath string

func init() {
	Cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the resource file to apply")
	Cmd.MarkFlagRequired("file")
}

var Cmd = &cobra.Command{
	Use:   "apply -f <file>",
	Short: "Apply Syncloud resources from a YAML file",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if _, err := os.Stat(filePath); err != nil {
			return err
		}

		parser := parser.NewResourceParser()
		validator := validatorsvc.NewValidator()
		builder := state.NewStateBuilder()

		dockerRecon := docker.NewDockerReconciler(components.Runner, components.Logger)
		k8sRecon := k8s.NewK8sReconciler(components.Runner, components.Logger)

		runtimes := []iApply.RuntimeReconciler{dockerRecon, k8sRecon}

		service := applyService.NewApplyService(
			components.Repo,
			builder,
			parser,
			validator,
			runtimes,
			components.Logger,
		)

		return service.Apply(filePath)
	},
}
