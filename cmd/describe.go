package cmd

import (
	"awasm-portfolio/internal/service"
	"fmt"

	"github.com/spf13/cobra"
)

func NewDescribeCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "describe [kind] [name] [namespace]",
		Short: "Describe a specific resource",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			kind, name, namespace := args[0], args[1], args[2]

			result, err := service.DescribeResource(kind, name, namespace)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Println(result)
		},
	}
}
