package get

import (
	resourses "synctl/cmd/get/nodes"
	nodes "synctl/cmd/get/resources"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "",
}

func init() {
	Cmd.AddCommand(resourses.Cmd)
	Cmd.AddCommand(nodes.Cmd)
}
