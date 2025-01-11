package cmd

import (
	"awasm-portfolio/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

func NewCreateCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "create [kind] [name]",
		Short: "Create a new resource",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			kind, name := args[0], args[1]
			namespace, _ := cmd.Flags().GetString("namespace")

			if namespace == "" {
				fmt.Println("Error: --namespace flag is required for create command")
				return
			}

			result, err := service.CreateResource(kind, name, namespace)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Println(result)
		},
	}
}
