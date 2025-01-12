package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewCreateCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:           "create [kind] [name]",
		Short:         "Create a new resource",
		SilenceErrors: false,
		Args:          cobra.ExactArgs(2),
		Example: `
# Create a profile in the default namespace
kubectl create profile john-doe

# Create a namespace
kubectl create namespace dev
		`,
		Run: func(cmd *cobra.Command, args []string) {
			kind, name := args[0], args[1]
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

			if allNamespaces {
				namespace = ""
			}

			result, err := service.CreateResource(kind, name, namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			cmd.Println(result)
		},
	}
}
