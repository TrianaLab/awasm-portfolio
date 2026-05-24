package cmd

import (
	"strings"

	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewGetCommand(repo *repository.InMemoryRepository) *cobra.Command {
	return &cobra.Command{
		Use:               "get [kind] [name]",
		Short:             "Get resources of a specific kind or a specific resource",
		Args:              cobra.RangeArgs(0, 2),
		ValidArgsFunction: completeResourceArgs(repo),
		Example: `
# Get all profiles in the default namespace
kubectl get profile

# Get a specific profile by name
kubectl get profile john-doe

# Get all resources of any kind across all namespaces
kubectl get all -A
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}

			flags := readFlags(cmd)
			lower := strings.ToLower(flags.output)
			if flags.output != "" && lower != "json" && lower != "yaml" {
				cmd.Printf("Error: wrong format '%s', valid ones are 'json' or 'yaml'\n", flags.output)
				return
			}

			kind := args[0]
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			result, err := service.Get(repo, kind, name, flags.namespace, flags.output)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			cmd.Println(result)
		},
	}
}
