package cli

import (
	"github.com/kolah/github-batch-updater/di"
	"github.com/spf13/cobra"
)

func RootCmd(app *di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use: "gbu",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(runCommand(app))

	return cmd
}
