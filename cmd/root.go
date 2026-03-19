package cmd

import (
	"synctl/cmd/apply"
	"synctl/cmd/daemon"
	"synctl/cmd/get"
	"synctl/cmd/inspect"
	"synctl/cmd/install"
	"synctl/cmd/version"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "synctl",
	Short: "synctl is a CLI tool for managing syncloud resources",
}

func init() {
	RootCmd.AddCommand(version.Cmd)
	RootCmd.AddCommand(inspect.Cmd)
	RootCmd.AddCommand(get.Cmd)
	RootCmd.AddCommand(install.Cmd)
	RootCmd.AddCommand(apply.Cmd)

	RootCmd.AddCommand(daemon.Cmd)
}
