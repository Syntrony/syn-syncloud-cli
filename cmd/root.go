package cmd

import (
	"github.com/Syntrony/syn-sycloud-cli/cmd/inspect"
	"github.com/Syntrony/syn-sycloud-cli/cmd/install"
	"github.com/Syntrony/syn-sycloud-cli/cmd/validate"
	"github.com/Syntrony/syn-sycloud-cli/cmd/version"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "synctl",
	Short: "synctl is a CLI tool for managing syncloud resources",
}

func init() {
	RootCmd.AddCommand(version.Cmd)
	RootCmd.AddCommand(inspect.Cmd)
	RootCmd.AddCommand(validate.Cmd)
	RootCmd.AddCommand(install.Cmd)
}
