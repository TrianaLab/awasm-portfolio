package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewGetCommand(svc *service.ResourceService) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "get [kind]",
		Short:         "List all resources of a given kind or all kinds",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := service.NormalizeResourceName(args[0]) // Normalize the kind

			// Special case for namespace
			if kind == "namespace" {
				resources, err := svc.ListResources(kind, "", true) // List namespaces without filtering
				if err != nil {
					return err
				}
				for _, res := range resources {
					cmd.Println(res.GetName())
				}
				return nil
			}

			// Handle other resources
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
			resources, err := svc.ListResources(kind, namespace, allNamespaces)
			if err != nil {
				return err
			}

			for _, res := range resources {
				cmd.Println(res.GetName())
			}
			return nil
		},
	}

	// Add namespace flags
	cmd.Flags().StringP("namespace", "n", "default", "Namespace to filter resources (default: 'default')")
	cmd.Flags().BoolP("all-namespaces", "A", false, "List resources across all namespaces")
	return cmd
}
