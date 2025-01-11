package cmd

import (
	"awasm-portfolio/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

func NewDescribeCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "describe [kind] [name]",
		Short: "Describe a specific resource",
		Args:  cobra.ExactArgs(2), // Require exactly 2 arguments: kind and name
		Run: func(cmd *cobra.Command, args []string) {
			kind, name := args[0], args[1]
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

			// Validate flags
			if allNamespaces {
				fmt.Println("Error: --all-namespaces is not supported for describe command")
				return
			}
			if namespace == "" {
				fmt.Println("Error: --namespace flag is required for describe command")
				return
			}

			result, err := service.DescribeResource(kind, name, namespace)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Println(result)
		},
	}
}
