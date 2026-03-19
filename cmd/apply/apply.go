package apply

import (
	"os"

	"synctl/internal/app"
	applyService "synctl/internal/application/services/apply"
	"synctl/internal/application/services/apply/parser"
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

		service := applyService.NewApplyService(
			components.Repo,
			builder,
			parser,
			validator,
			components.Logger,
		)

		return service.Apply(filePath)
	},
}
