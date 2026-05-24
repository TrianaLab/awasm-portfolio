package cmd

import (
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDeleteCommand(repo *repository.InMemoryRepository) *cobra.Command {
	return &cobra.Command{
		Use:               "delete [kind] [name]",
		Short:             "Delete a resource",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: completeResourceArgs(repo),
		Example: `
# Delete the profile john-doe in the dev namespace
kubectl delete profile john-doe -n dev

# Delete a namespace
kubectl delete namespace dev
		`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := readFlags(cmd)
			if flags.output != "" {
				cmd.Println("Error: output flag is not supported in this command")
				return
			}

			result, err := service.Delete(repo, args[0], args[1], flags.namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			cmd.Println(result)
		},
	}
}
