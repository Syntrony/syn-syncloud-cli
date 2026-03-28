package delete

import (
	"synctl/internal/app"
	commonParser "synctl/internal/application/services/common/parser"
	commonValidator "synctl/internal/application/services/common/validator"
	"synctl/internal/application/services/delete/mutation"
	mutationSvc "synctl/internal/application/services/mutation"
	"synctl/internal/application/services/state"

	"github.com/spf13/cobra"
)

var (
	name     string
	id       string
	filePath string
)

func init() {
	Cmd.Flags().StringVarP(&filePath, "file", "f", "", "delete by resource file")
	Cmd.MarkFlagRequired("file")
}

var Cmd = &cobra.Command{
	Use:   "delete -f <file>",
	Short: "Delete Syncloud resources from a YAML file",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if name == "" && id == "" && filePath == "" {
			return cmd.Help()
		}

		parser := commonParser.NewResourceParser()
		validator := commonValidator.NewValidator()
		builder := state.NewStateBuilder()

		deleteMutation := mutation.NewDeleteMutation(components.Repo)
		mutationService := mutationSvc.NewMutationService(
			parser,
			validator,
			components.Repo,
			builder,
			components.Logger,
		)

		components.Logger.Info("Delete resources...")

		if filePath != "" {
			return mutationService.Execute(filePath, deleteMutation)
		}

		components.Logger.Info("Delete by name/id not yet implemented")
		return nil
	},
}
