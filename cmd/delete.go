package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(svc *service.ResourceService) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "delete [kind] [name]",
		Short:         "Delete a specific resource and its child resources",
		Args:          cobra.MinimumNArgs(2),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := service.NormalizeResourceName(args[0]) // Normalize the kind
			name := args[1]

			// Special case for namespace
			if kind == "namespace" {
				message, err := svc.DeleteResourceWithNamespace(kind, name, "")
				if err != nil {
					return err
				}
				cmd.Println(message)
				return nil
			}

			// Handle other resources
			namespace, _ := cmd.Flags().GetString("namespace")
			message, err := svc.DeleteResourceWithNamespace(kind, name, namespace)
			if err != nil {
				return err
			}

			cmd.Println(message)
			return nil
		},
	}

	// Add namespace flag
	cmd.Flags().StringP("namespace", "n", "default", "Namespace of the resource (default: 'default')")
	return cmd
}
