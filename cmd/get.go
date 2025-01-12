package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewGetCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "get [kind] [name]",
		Short: "Get resources of a specific kind or a specific resource",
		Args:  cobra.RangeArgs(1, 2), // Accept 1 or 2 arguments
		Run: func(cmd *cobra.Command, args []string) {
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
