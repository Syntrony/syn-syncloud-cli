package validate

import (
	"fmt"

	"github.com/Syntrony/syn-sycloud-cli/internal/validator"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate system requirements for Syncloud",
	Run: func(cmd *cobra.Command, args []string) {

		result := validator.RunValidator()

		fmt.Println("Validating environment... \n")

		if result.DockerInstalled {
			fmt.Println("docker is installed")
		} else {
			fmt.Println("docker is not installed")
		}

		if result.KubectlInstalled {
			fmt.Println("kubectl is installed")
		} else {
			fmt.Println("kubectl is not installed")
		}
	},
}
