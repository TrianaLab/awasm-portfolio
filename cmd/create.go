package cmd

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewCreateCommand(svc *service.ResourceService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [kind] [name]",
		Short: "Create a new resource",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := service.NormalizeResourceName(args[0]) // Normalize the kind
			name := args[1]
			namespace, _ := cmd.Flags().GetString("namespace")

			// Use ResourceFactory to create the resource
			resourceFactory := factory.ResourceFactory{}
			resource := resourceFactory.Create(kind, map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			})

			return svc.CreateResource(kind, resource)
		},
	}

	// Add namespace flag
	cmd.Flags().StringP("namespace", "n", "default", "Namespace for the resource (default: 'default')")
	return cmd
}
