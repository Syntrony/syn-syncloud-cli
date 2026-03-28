package apply

import (
	"os"

	"synctl/internal/app"
	mutation "synctl/internal/application/services/apply"
	commonParser "synctl/internal/application/services/common/parser"
	commonValidator "synctl/internal/application/services/common/validator"
	mutationSvc "synctl/internal/application/services/mutation"
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

		parser := commonParser.NewResourceParser()
		validator := commonValidator.NewValidator()
		builder := state.NewStateBuilder()

		applyMutation := mutation.NewApplyMutation(components.Repo)
		mutationService := mutationSvc.NewMutationService(
			parser,
			validator,
			components.Repo,
			builder,
			components.Logger,
		)

		components.Logger.Info("Apply resources...")

		if filePath != "" {
			return mutationService.Execute(filePath, applyMutation)
		}

		components.Logger.Info("Apply single resource not yet implemented")
		return nil
	},
}
