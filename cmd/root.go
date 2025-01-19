package cmd

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/service"
	"awasm-portfolio/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCommand(repo *repository.InMemoryRepository) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:           "kubectl",
		Short:         "A CLI tool for managing resources",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	resourceService := service.NewResourceService(repo, rootCmd)

	// Remove the default 'help' subcommand entirely
	rootCmd.SetHelpCommand(nil)

	// Set the custom help function for --help flag usage
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Println("A CLI tool for managing resources")
		cmd.Println("\nUsage:")
		cmd.Println("  kubectl [command] [resource]")
		cmd.Println("\nAliases:")
		cmd.Println("  kubectl, k")
		cmd.Println("\nAvailable Commands:")
		cmd.Println("  create      Create a new resource")
		cmd.Println("  delete      Delete a resource")
		cmd.Println("  describe    Describe a specific resource")
		cmd.Println("  get         Get resources of a specific kind or a specific resource")
		cmd.Println("  version     Get the current version of the application")

		cmd.Println("\nAvailable Resources:")
		resources := util.SupportedResources()
		canonical := map[string]string{}
		aliases := map[string][]string{}

		for alias, kind := range resources {
			if _, exists := canonical[kind]; !exists {
				canonical[kind] = alias
			} else {
				aliases[kind] = append(aliases[kind], alias)
			}
		}
		for kind, aliasList := range aliases {
			cmd.Printf("  %s  (aliases: %s)\n", kind, fmt.Sprintf("%s", aliasList))
		}

		cmd.Println("\nFlags:")
		cmd.Println("  -A, --all-namespaces     Include resources from all namespaces")
		cmd.Println("  -h, --help               help for kubectl")
		cmd.Println("  -n, --namespace string   Specify the namespace (default \"default\")")
		cmd.Println("\nUse \"kubectl [command] --help\" for more information about a command.")
	})

	// Add subcommands
	rootCmd.AddCommand(NewCreateCommand(resourceService))
	rootCmd.AddCommand(NewDeleteCommand(resourceService))
	rootCmd.AddCommand(NewGetCommand(resourceService))
	rootCmd.AddCommand(NewDescribeCommand(resourceService))
	rootCmd.AddCommand(NewVersionCommand())

	// Add global persistent flags
	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Specify the namespace")
	rootCmd.PersistentFlags().BoolP("all-namespaces", "A", false, "Include resources from all namespaces")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Specify the output format, either 'json' or 'yaml'")

	// Disable autocompletion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}
