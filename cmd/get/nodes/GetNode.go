package nodes

import (
	"synctl/internal/app"
	outputs "synctl/internal/application/outputs"
	services "synctl/internal/application/services/get"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "nodes",
	Short: "List existing nodes in the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		nService := services.GetNodeService{
			Repo: components.Repo,
		}

		result, err := nService.Execute()
		if err != nil {
			return err
		}

		if len(result) == 0 {
			components.Logger.Info("No nodes found in system!!!")
			return nil
		}

		outputs.PrintNodes(result)
		return nil

	},
}
