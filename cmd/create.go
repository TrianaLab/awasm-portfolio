package cmd

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewCreateCommand(svc *service.ResourceService) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create [kind] [name]",
		Short:         "Create a new resource",
		Args:          cobra.MinimumNArgs(2),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := service.NormalizeResourceName(args[0]) // Normalize the kind
			name := args[1]

			// Special case for namespace
			if kind == "namespace" {
				resource := &types.Namespace{Name: name}
				message, err := svc.CreateResource(kind, resource)
				if err != nil {
					return err
				}
				cmd.Println(message)
				return nil
			}

			// Handle other resources
			namespace, _ := cmd.Flags().GetString("namespace")
			resourceFactory := factory.ResourceFactory{}
			resource := resourceFactory.Create(kind, map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			})

			message, err := svc.CreateResource(kind, resource)
			if err != nil {
				return err
			}
			cmd.Println(message)
			return nil
		},
	}

	// Add namespace flag
	cmd.Flags().StringP("namespace", "n", "default", "Namespace for the resource (default: 'default')")
	return cmd
}
