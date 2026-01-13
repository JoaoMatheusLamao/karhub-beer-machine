package root

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root CLI command.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "karhub-cli",
		Short: "Karhub Beer Machine CLI",
		Long:  "CLI for administrative and maintenance tasks of the Karhub Beer Machine",
	}

	cmd.AddCommand(
		newSeedCommand(),
	)

	return cmd
}

// Execute runs the root command.
func Execute() error {
	return NewRootCommand().Execute()
}
