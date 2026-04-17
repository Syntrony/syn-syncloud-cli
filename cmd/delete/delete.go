package delete

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"synctl/internal/app"
	commonParser "synctl/internal/application/services/common/parser"
	commonValidator "synctl/internal/application/services/common/validator"
	"synctl/internal/application/services/daemon/reconciler/docker"
	"synctl/internal/application/services/daemon/reconciler/k8s"
	mutation "synctl/internal/application/services/delete"
	getSvc "synctl/internal/application/services/get"
	mutationSvc "synctl/internal/application/services/mutation"
	"synctl/internal/application/services/state"
	"synctl/internal/domain"
	filters "synctl/internal/domain/filters"

	"github.com/spf13/cobra"
)

func isPathSafe(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal not allowed: %s", filePath)
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path must be a file, not a directory: %s", filePath)
	}

	return nil
}

var (
	name     string
	id       string
	filePath string
	dryRun   bool
)

func init() {
	Cmd.Flags().StringVarP(&filePath, "file", "f", "", "delete by resource file")
	Cmd.Flags().StringVarP(&name, "name", "n", "", "delete by resource name")
	Cmd.Flags().StringVarP(&id, "id", "I", "", "delete by resource id")
	Cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview deletions without applying them")
}

var Cmd = &cobra.Command{
	Use:   "delete -f <file> | -n <name> | -I <id>",
	Short: "Delete Syncloud resources by file, name, or ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if name == "" && id == "" && filePath == "" {
			return fmt.Errorf("requires at least one flag: -f, -n, or -I")
		}

		parser := commonParser.NewResourceParser()
		validator := commonValidator.NewValidator()
		builder := state.NewStateBuilder()

		deleteMutation := mutation.NewDeleteMutation(components.Repo)
		mutationService := mutationSvc.NewMutationService(
			parser,
			validator,
			components.Repo,
			builder,
			components.Logger,
		).WithInspector(components.Inspector)

		if dryRun {
			mutationService.AsDryRun()
		}

		components.Logger.Info("Delete resources...")

		var resources []*domain.Resource

		if filePath != "" {
			if err := isPathSafe(filePath); err != nil {
				return err
			}
			parsed, err := parser.Parse(filePath)
			if err != nil {
				return err
			}
			resources = parsed
		} else {
			filter := filters.GetResourceFilter{
				Name: name,
				Id:   id,
			}
			getService := getSvc.GetResourceService{
				Repo: components.Repo,
			}
			result, err := getService.GetResourceBy(filter)
			if err != nil {
				return err
			}
			resources = result
		}

		if len(resources) == 0 {
			components.Logger.Info("No resources found to delete")
			return nil
		}

		if err := mutationService.ExecuteResource(resources, deleteMutation); err != nil {
			return err
		}

		if dryRun {
			return nil
		}

		dockerRecon := docker.NewDockerReconciler(components.Runner, components.Logger)
		k8sRecon := k8s.NewK8sReconciler(components.Runner, components.Logger)

		for _, rt := range []interface {
			Runtime() string
			Reconcile(domain.ReconcileContext) error
		}{dockerRecon, k8sRecon} {
			var rtResources []*domain.Resource
			var rtActions []*domain.Action

			for _, r := range resources {
				if r.Runtime == rt.Runtime() {
					rtResources = append(rtResources, r)
					rtActions = append(rtActions, &domain.Action{
						Type:     domain.ActionDelete,
						Resource: r,
					})
				}
			}

			if len(rtResources) == 0 {
				continue
			}

			subCtx := domain.ReconcileContext{
				Resources: rtResources,
				Actions:   rtActions,
			}

			if err := rt.Reconcile(subCtx); err != nil {
				components.Logger.Error(fmt.Sprintf("infra reconcile error: %s", err))
			}
		}

		return nil
	},
}
