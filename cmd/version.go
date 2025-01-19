package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "dev"

// NewVersionCommand returns a new cobra command for the version.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  "Display the current application version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "App version: %s\n", appVersion)
		},
	}
}
