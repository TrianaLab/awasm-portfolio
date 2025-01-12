package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewGetCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "get [kind] [name]",
		Short: "Get resources of a specific kind or a specific resource",
		Args:  cobra.RangeArgs(0, 2),
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
				cmd.Help()
				return
			}

			kind := args[0]
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

			if allNamespaces {
				namespace = ""
			}

			result, err := service.GetResources(kind, name, namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			cmd.Println(result)
		},
	}
}
