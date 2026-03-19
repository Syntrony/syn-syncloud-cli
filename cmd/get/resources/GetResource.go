package resources

import (
	"synctl/internal/app"
	outputs "synctl/internal/application/outputs"
	services "synctl/internal/application/services/get"
	filters "synctl/internal/domain/filters"

	"github.com/spf13/cobra"
)

var (
	name    string
	runtime string
	nodeId  string
	kind    string
	id      string
)

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n", "", "Filter by resource name")
	Cmd.Flags().StringVarP(&kind, "kind", "k", "", "Filter by resource kind")
	Cmd.Flags().StringVarP(&runtime, "runtime", "r", "", "Filter by resource runtime")
	Cmd.Flags().StringVarP(&id, "id", "i", "", "Filter by resource id")
}

var Cmd = &cobra.Command{
	Use:   "resources",
	Short: "List existing resources in every node",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		rService := services.GetResourceService{
			Repo: components.Repo,
		}

		filter := filters.GetResourceFilter{
			Name:    name,
			Runtime: runtime,
			NodeId:  nodeId,
			Kind:    kind,
			Id:      id,
		}

		result, err := rService.Execute(filter)

		if err != nil {
			return err
		}

		if len(result) == 0 {
			components.Logger.Info("No resources found in system!!!")
			return nil
		}

		outputs.PrintResources(result)

		return nil
	},
}
