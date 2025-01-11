package cmd

import (
	"awasm-portfolio/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [kind] [name]",
		Short: "Delete a resource",
		Args:  cobra.ExactArgs(2), // Require exactly 2 arguments: kind and name
		Run: func(cmd *cobra.Command, args []string) {
			kind, name := args[0], args[1]
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

			// Validate flags
			if allNamespaces {
				fmt.Println("Error: --all-namespaces is not supported for delete command")
				return
			}
			if namespace == "" {
				fmt.Println("Error: --namespace flag is required for delete command")
				return
			}

			result, err := service.DeleteResource(kind, name, namespace)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			cmd.Println(result)
		},
	}
}
