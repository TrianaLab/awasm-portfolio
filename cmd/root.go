package cmd

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewRootCommand(repo *repository.InMemoryRepository) *cobra.Command {
	resourceService := service.NewResourceService(repo)

	rootCmd := &cobra.Command{
		Use:           "kubectl",
		Aliases:       []string{"k"},
		Short:         "A CLI tool for managing resources",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(NewCreateCommand(resourceService))
	rootCmd.AddCommand(NewDeleteCommand(resourceService))
	rootCmd.AddCommand(NewGetCommand(resourceService))
	rootCmd.AddCommand(NewDescribeCommand(resourceService))

	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Specify the namespace")
	rootCmd.PersistentFlags().BoolP("all-namespaces", "A", false, "Include resources from all namespaces")

	return rootCmd
}
