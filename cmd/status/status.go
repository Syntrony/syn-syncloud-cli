package status

import (
	"fmt"
	"synctl/internal/app"
	statusservice "synctl/internal/application/services/status"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "Show platform health status",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		svc := statusservice.NewStatusService(components.Repo, components.Runner)
		s, err := svc.Execute()
		if err != nil {
			return err
		}

		fmt.Println("╔══════════════════════════════════════╗")
		fmt.Println("║      Syncloud Platform - Status      ║")
		fmt.Println("╚══════════════════════════════════════╝")
		fmt.Println()

		fmt.Printf("  Docker:      %s\n", s.Docker)
		fmt.Printf("  Kubernetes:  %s\n", s.Kubernetes)
		fmt.Println()

		if s.StateFound {
			fmt.Printf("  Platform:    installed (v%s)\n", s.Version)
			fmt.Printf("  Cluster:     %s [%s]\n", s.ClusterName, s.ClusterMode)
			fmt.Printf("  Nodes:       %d\n", s.NodeCount)
			fmt.Printf("  Resources:   %d\n", s.ResourceCount)
		} else {
			fmt.Println("  Platform:    not installed (run: synctl install)")
		}
		fmt.Println()

		return nil
	},
}
