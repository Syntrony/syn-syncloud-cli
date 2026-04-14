package logs

import (
	"fmt"
	"strings"
	"synctl/internal/app"
	getservice "synctl/internal/application/services/get"
	logsservice "synctl/internal/application/services/logs"
	"synctl/internal/domain/filters"

	"github.com/spf13/cobra"
)

var (
	flagName      string
	flagRuntime   string
	flagNamespace string
	flagFollow    bool
	flagTail      int
)

var Cmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs of a resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flagName == "" {
			return fmt.Errorf("resource name is required (-n)")
		}

		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		runtime := flagRuntime
		namespace := flagNamespace

		// Resolve runtime from state if not provided
		if runtime == "" {
			getSvc := &getservice.GetResourceService{Repo: components.Repo}
			results, err := getSvc.GetResourceBy(filters.GetResourceFilter{Name: flagName})
			if err != nil {
				return err
			}
			if len(results) > 0 {
				runtime = results[0].Runtime
				if ns, ok := results[0].Spec["namespace"].(string); ok && namespace == "" {
					namespace = ns
				}
			}
		}

		if runtime == "" {
			return fmt.Errorf("could not determine runtime for '%s': use -r to specify (docker|kubernetes)", flagName)
		}

		svc := logsservice.NewLogsService(components.Runner)

		output, err := svc.Fetch(logsservice.LogsOptions{
			Name:      flagName,
			Runtime:   strings.ToLower(runtime),
			Namespace: namespace,
			Follow:    flagFollow,
			Tail:      flagTail,
		})
		if err != nil {
			return err
		}

		fmt.Print(output)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVarP(&flagName, "name", "n", "", "Resource name")
	Cmd.Flags().StringVarP(&flagRuntime, "runtime", "r", "", "Runtime (docker|kubernetes)")
	Cmd.Flags().StringVar(&flagNamespace, "namespace", "", "Kubernetes namespace")
	Cmd.Flags().BoolVarP(&flagFollow, "follow", "f", false, "Follow log output")
	Cmd.Flags().IntVar(&flagTail, "tail", 0, "Number of last lines to show (0 = all)")
}
