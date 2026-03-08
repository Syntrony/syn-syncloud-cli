package nodes

import (
	outputs "synctl/internal/application/outputs"
	services "synctl/internal/application/services/get"
	"synctl/internal/domain/persistence"
	loggerinfra "synctl/internal/logger"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "nodes",
	Short: "List existing nodes in the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := &persistence.StateRepository{
			Path: "helpers/syncloud-state.json",
		}

		logger := loggerinfra.NewConsoleLogger()

		nService := services.GetNodeService{
			Repo: repo,
		}

		result, err := nService.Execute()
		if err != nil {
			return err
		}

		if len(result) == 0 {
			logger.Info("No nodes found in system!!!")
			return nil
		}

		outputs.PrintNodes(result)
		return nil

	},
}
