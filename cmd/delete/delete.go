package delete

import (
	"fmt"
	"synctl/internal/app"
	commonParser "synctl/internal/application/services/common/parser"
	commonValidator "synctl/internal/application/services/common/validator"
	mutation "synctl/internal/application/services/delete"
	getSvc "synctl/internal/application/services/get"
	mutationSvc "synctl/internal/application/services/mutation"
	"synctl/internal/application/services/state"
	filters "synctl/internal/domain/filters"

	"github.com/spf13/cobra"
)

var (
	name     string
	id       string
	filePath string
)

func init() {
	Cmd.Flags().StringVarP(&filePath, "file", "f", "", "delete by resource file")
	Cmd.Flags().StringVarP(&name, "name", "n", "", "delete by resource name")
	Cmd.Flags().StringVarP(&id, "id", "I", "", "delete by resource id")

}

var Cmd = &cobra.Command{
	Use:   "delete -f <file> | -n <name> | -I <id>",
	Short: "Delete Syncloud resources by file, name, or ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if name == "" && id == "" && filePath == "" {
			return fmt.Errorf("requires at least one flag: -f, -n, or -I")
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
			return mutationService.ExecuteFile(filePath, deleteMutation)
		}

		filter := filters.GetResourceFilter{
			Name: name,
			Id:   id,
		}

		getService := getSvc.GetResourceService{
			Repo: components.Repo,
		}

		result, err := getService.GetResourceBy(filter)

		if err != nil {
			return err
		}

		return mutationService.ExecuteResource(result, deleteMutation)
	},
}
