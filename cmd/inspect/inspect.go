package inspect

import (
	services "synctl/internal/application/services"
	"synctl/internal/domain/persistence"
	loggerinfra "synctl/internal/logger"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect existing resources on the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := &persistence.StateRepository{
			Path: "helpers/syncloud-state.json",
		}

		logger := loggerinfra.NewConsoleLogger()

		service := services.InspectService{
			Repo: repo,
		}

		result, err := service.Execute()
		if err != nil {
			return err
		}

		if !result.Found {
			logger.Info("Syncloud not initialized")
			return nil
		}

		logger.Info("Syncloud State: FOUND")
		logger.Info("Version: " + result.Version)
		logger.Info("Mode: " + result.Mode)
		logger.Info("Nodes: " + string(rune(result.NodeCount+'0')))
		logger.Info("Resources: " + string(rune(result.ResCount+'0')))

		return nil

	},
}
