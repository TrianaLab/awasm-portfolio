package cmd

import (
	"awasm-portfolio/internal/service"
	"fmt"

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
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

			if kind == "all" {
				// Handle "all" special case
				resourcesByKind, err := svc.ListAllResources(namespace, allNamespaces)
				if err != nil {
					return err
				}

				// Print grouped resources
				for kind, resources := range resourcesByKind {
					cmd.Printf("Kind: %s\n", kind)
					for _, res := range resources {
						cmd.Printf("  %s\n", res.GetName())
					}
					cmd.Println()
				}
				return nil
			}

			// Validate specific kind
			if !service.IsValidResource(kind) {
				return fmt.Errorf("error: unknown resource type '%s'", kind)
			}

			// Handle specific kind
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
