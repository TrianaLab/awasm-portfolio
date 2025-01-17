package cmd

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/service"
	"awasm-portfolio/internal/util"
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCommand(repo *repository.InMemoryRepository) *cobra.Command {
	resourceService := service.NewResourceService(repo)

	rootCmd := &cobra.Command{
		Use:           "kubectl",
		Short:         "A CLI tool for managing resources",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Define the custom help function
	customHelpFunc := func(cmd *cobra.Command, args []string) {
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
		cmd.Println("  help        Help about any command")

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
	}

	// Set the custom help function
	rootCmd.SetHelpFunc(customHelpFunc)
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help",
		Short: "Help about any command",
		Run: func(cmd *cobra.Command, args []string) {
			customHelpFunc(rootCmd, args)
		},
	})

	// Add subcommands
	rootCmd.AddCommand(NewCreateCommand(resourceService))
	rootCmd.AddCommand(NewDeleteCommand(resourceService))
	rootCmd.AddCommand(NewGetCommand(resourceService))
	rootCmd.AddCommand(NewDescribeCommand(resourceService))

	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Specify the namespace")
	rootCmd.PersistentFlags().BoolP("all-namespaces", "A", false, "Include resources from all namespaces")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd
}
