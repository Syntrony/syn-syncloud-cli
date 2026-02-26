package nodes

import (
	"fmt"
	outputs "synctl/internal/application/outputs"
	repository "synctl/internal/application/repository"
	services "synctl/internal/application/services/get"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "nodes",
	Short: "List existing nodes in the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.StateRepository()

		nService := services.GetNodeService{
			Repo: repo,
		}

		result, err := nService.Execute()
		if err != nil {
			return err
		}

		if len(result) == 0 {
			fmt.Println("No nodes found in system!!!")
			return nil
		}

		outputs.PrintNodes(result)
		return nil

	},
}
