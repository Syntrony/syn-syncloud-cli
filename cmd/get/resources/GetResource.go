package resources

import (
	"fmt"
	"synctl/internal/application/outputs"
	repository "synctl/internal/application/repository"
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
	// Cmd.Flags().StringVarP(&nodeId, "node-id", "nid", "", "Filter by node id")
	Cmd.Flags().StringVarP(&id, "id", "i", "", "Filter by resource id")
}

var Cmd = &cobra.Command{
	Use:   "resources",
	Short: "List existing resources in every node",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.StateRepository()

		rService := services.GetResourceService{
			Repo: repo,
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
			fmt.Println("No resources found in system!!!")

			return nil
		}

		outputs.PrintResources(result)

		return nil
	},
}
