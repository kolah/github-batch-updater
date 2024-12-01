package cli

import (
	"github.com/kolah/github-batch-updater/internal/di"

	"github.com/spf13/cobra"
)

func runCommand(app *di.Container) *cobra.Command {
	dryRun := false
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the application",
		Args:  cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			c, err := app.LoadBatch(fileName)
			if err != nil {
				return err
			}

			app.BatchProcessingService().Process(cmd.Context(), c.ToInput(dryRun))

			return nil
		},
	}
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "dry run")

	return cmd
}
