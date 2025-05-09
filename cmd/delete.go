package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [kind] [name]",
		Short: "Delete a resource",
		Args:  cobra.ExactArgs(2),
		Example: `
# Delete the profile john-doe in the dev namespace
kubectl delete profile john-doe -n dev

# Delete a namespace
kubectl delete namespace dev
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

			result, err := service.DeleteResource(kind, name, namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			cmd.Println(result)
		},
	}
}
