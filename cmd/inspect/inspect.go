package inspect

import (
	"fmt"
	services "synctl/internal/application/services"
	state "synctl/internal/state"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect existing resources on the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := &state.StateRepository{
			Path: "helpers/state.json",
		}

		service := services.InspectService{
			Repo: repo,
		}

		result, err := service.Execute()
		if err != nil {
			return err
		}

		if !result.Found {
			fmt.Println("Syncloud not initialized")
			return nil
		}

		fmt.Println("Syncloud State: FOUND")
		fmt.Println("Version:", result.Version)
		fmt.Println("Mode:", result.Mode)
		fmt.Println("Nodes:", result.NodeCount)
		fmt.Println("Resources:", result.ResCount)

		return nil

	},
}
