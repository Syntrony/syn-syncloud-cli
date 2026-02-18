package install

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "install",
	Short: "Install Syncloud Platform",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Syncloud installation...")
	},
}
