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
# Create a namespace
kubectl create namespace dev

# Create a profile in the dev namespace
kubectl create profile john-doe -n dev
		`,
		Run: func(cmd *cobra.Command, args []string) {
			kind, name := args[0], args[1]
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
			formatOutput, _ := cmd.Flags().GetString("output")

			if formatOutput != "" {
				cmd.Println("Error: output flag is not supported in this command")
				return
			}

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
