package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(svc *service.ResourceService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [kind] [name]",
		Short: "Delete a specific resource and its child resources",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := args[0]
			name := args[1]
			namespace, _ := cmd.Flags().GetString("namespace")

			return svc.DeleteResourceWithNamespace(kind, name, namespace)
		},
	}

	// Add namespace flag
	cmd.Flags().StringP("namespace", "n", "default", "Namespace of the resource")
	return cmd
}
