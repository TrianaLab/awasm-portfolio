package cmd

import (
	"awasm-portfolio/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

func NewCreateCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "create [kind] [name] [namespace]",
		Short: "Create a new resource",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			kind, name, namespace := args[0], args[1], args[2]

			result, err := service.CreateResource(kind, name, namespace)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Println(result)
		},
	}
}
