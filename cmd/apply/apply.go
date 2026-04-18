package apply

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"synctl/internal/app"
	mutation "synctl/internal/application/services/apply"
	commonParser "synctl/internal/application/services/common/parser"
	commonValidator "synctl/internal/application/services/common/validator"
	mutationSvc "synctl/internal/application/services/mutation"
	"synctl/internal/application/services/state"

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
	filePath string
	dryRun   bool
)

func init() {
	Cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the resource file to apply")
	Cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")
	Cmd.MarkFlagRequired("file")
}

var Cmd = &cobra.Command{
	Use:   "apply -f <file>",
	Short: "Apply Syncloud resources from a YAML file",
	RunE: func(cmd *cobra.Command, args []string) error {

		components := app.NewComponents()
		components.Init()
		components.WithDefaultRepo()

		if err := isPathSafe(filePath); err != nil {
			return err
		}

		parser := commonParser.NewResourceParser()
		validator := commonValidator.NewValidator()
		builder := state.NewStateBuilder()

		applyMutation := mutation.NewApplyMutation(components.Repo)
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

		components.Logger.Info("Apply resources...")

		if filePath != "" {
			if err := mutationService.ExecuteFile(filePath, applyMutation); err != nil {
				return err
			}
			if !dryRun {
				components.Logger.Info("Triggering background reconcile...")
				app.TriggerReconcile()
			}
			return nil
		}

		components.Logger.Info("Apply single resource not yet implemented")
		return nil
	},
}
