package inspect

import (
	"fmt"

	"github.com/Syntrony/syn-sycloud-cli/internal/system"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect existing resources on the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Inspecting resources... \n")

		if !system.CommandExists("docker") {
			fmt.Println("Docker is not installed.")
		} else {
			outDocker, err := system.RunCommand("docker", "ps")

			if err != nil {
				fmt.Println("Docker found but failed", err)
			} else {
				fmt.Println("Existing docker containers:")
				fmt.Println(outDocker)
			}
		}

		fmt.Println()

		if !system.CommandExists("kubectl") {
			fmt.Println("kubectl is not installed.")
		} else {
			outK8s, err := system.RunCommand("kubectl", "get", "pods", "--all-namespaces")
			if err != nil {
				fmt.Println("kubectl found but failed:", err)
			} else {
				fmt.Println("Existing Kubernetes pods:")
				fmt.Println(outK8s)
			}
		}
	},
}
