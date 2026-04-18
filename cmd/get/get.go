package get

import (
	dns "synctl/cmd/get/dns"
	nodes "synctl/cmd/get/nodes"
	resources "synctl/cmd/get/resources"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "List specific resources / nobes given filters",
}

func init() {
	Cmd.AddCommand(resources.Cmd)
	Cmd.AddCommand(nodes.Cmd)
	Cmd.AddCommand(dns.Cmd)
}
