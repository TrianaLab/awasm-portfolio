package cmd

import (
	"awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(svc *service.ResourceService) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "delete [kind] [name]",
		Short:         "Delete a specific resource and its child resources",
		Args:          cobra.ExactArgs(2),
		SilenceUsage:  true, // Prevent help text on error
		SilenceErrors: true, // Prevent error stack on error
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := service.NormalizeResourceName(args[0]) // Normalize the kind
			name := args[1]
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
	cmd.Flags().StringP("namespace", "n", "default", "Namespace of the resource")
	return cmd
}
