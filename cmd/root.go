package cmd

import (
	"synctl/cmd/apply"
	"synctl/cmd/daemon"
	"synctl/cmd/delete"
	"synctl/cmd/get"
	"synctl/cmd/inspect"
	"synctl/cmd/install"
	"synctl/cmd/version"
	"synctl/internal/logger"

	"github.com/spf13/cobra"
)

var log *logger.Logger

func preRun(cmd *cobra.Command, args []string) {
	log = logger.NewConsoleLogger()
	log.Command(">>> Running: synctl " + cmd.Name())
}

func postRun(cmd *cobra.Command, args []string) {
	log.Command("<<< Completed: synctl " + cmd.Name())
	log = nil
}

var RootCmd = &cobra.Command{
	Use:   "synctl",
	Short: "synctl is a CLI tool for managing syncloud resources",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		preRun(cmd, args)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		postRun(cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(version.Cmd)
	RootCmd.AddCommand(inspect.Cmd)
	RootCmd.AddCommand(get.Cmd)
	RootCmd.AddCommand(install.Cmd)
	RootCmd.AddCommand(apply.Cmd)
	RootCmd.AddCommand(delete.Cmd)

	RootCmd.AddCommand(daemon.Cmd)
}
