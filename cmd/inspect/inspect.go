package inspect

import (
	"fmt"
	"synctl/internal/app"
	services "synctl/internal/application/services"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect existing resources on the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		service := services.InspectService{
			Repo: components.Repo,
		}

		result, err := service.Execute()
		if err != nil {
			return err
		}

		if !result.Found {
			components.Logger.Info("Syncloud not initialized")
			return nil
		}

		components.Logger.Info("Syncloud State: FOUND")
		components.Logger.Info("Version: " + result.Version)
		components.Logger.Info("Mode: " + result.Mode)
		components.Logger.Info(fmt.Sprintf("Nodes: %d", result.NodeCount))
		components.Logger.Info(fmt.Sprintf("Resources: %d", result.ResCount))

		return nil
	},
}
