package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDescribeCommand(service service.ResourceService) *cobra.Command {
	return &cobra.Command{
		Use:   "describe [kind] [name]",
		Short: "Describe a specific resource",
		Args:  cobra.RangeArgs(0, 2),
		Example: `
# Describe the profile jane-doe in the dev namespace
kubectl describe profile jane-doe -n dev
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
			formatOutput, _ := cmd.Flags().GetString("output")

			if formatOutput != "" {
				cmd.Println("Error: output flag is not supported in this command")
				return
			}

			if allNamespaces {
				namespace = ""
			}

			result, err := service.DescribeResource(kind, name, namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}

			cmd.Println(result)
		},
	}
}
