package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewRootCommand(svc *service.ResourceService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kubectl",
		Short: "A CLI tool for managing resources",
	}

	// Add subcommands
	rootCmd.AddCommand(NewCreateCommand(svc))
	rootCmd.AddCommand(NewDeleteCommand(svc))
	rootCmd.AddCommand(NewGetCommand(svc))
	rootCmd.AddCommand(NewDescribeCommand(svc))

	return rootCmd
}
