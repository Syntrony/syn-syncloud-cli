package rollback

import (
	"fmt"
	"os"

	"synctl/internal/app"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "rollback",
	Short: "Restore platform state from the last backup (.bak)",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		state, err := components.Repo.Load()
		if err != nil {
			return fmt.Errorf("could not load current state: %w", err)
		}

		_ = state // current state loaded for reference

		statePath := app.VPSStatePath
		if components.EnvDetector.IsContainer() {
			statePath = app.ContainerStatePath
		}

		backupPath := statePath + ".bak"

		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			return fmt.Errorf("no backup found at %s", backupPath)
		}

		backup, err := os.ReadFile(backupPath)
		if err != nil {
			return fmt.Errorf("failed to read backup: %w", err)
		}

		components.Logger.Info(fmt.Sprintf("Restoring state from %s...", backupPath))

		if err := os.WriteFile(statePath, backup, 0600); err != nil {
			return fmt.Errorf("failed to restore backup: %w", err)
		}

		components.Logger.Success("State rolled back successfully. Run 'synctl reconcile' to realign runtime.")
		return nil
	},
}
