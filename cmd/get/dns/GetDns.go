package dns

import (
	"synctl/internal/app"
	outputs "synctl/internal/application/outputs"
	services "synctl/internal/application/services/get"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "dns",
	Short: "List DNS records in the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		dService := services.GetDnsService{
			Repo: components.Repo,
		}

		result, err := dService.Execute()
		if err != nil {
			return err
		}

		if len(result) == 0 {
			components.Logger.Info("No DNS records found in system!!!")
			return nil
		}

		outputs.PrintDnsRecords(result)
		return nil
	},
}
