package version

import (
	"synctl/internal/logger"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of synctl",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.NewConsoleLogger()
		logger.Info("synctl version " + version)
	},
}
