package cmd

import (
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewCreateCommand(repo *repository.InMemoryRepository) *cobra.Command {
	return &cobra.Command{
		Use:           "create [kind] [name]",
		Short:         "Create a new resource",
		SilenceErrors: false,
		Args:          cobra.ExactArgs(2),
		Example: `
# Create a namespace
kubectl create namespace dev

# Create a profile in the dev namespace
kubectl create profile john-doe -n dev
		`,
		Run: func(cmd *cobra.Command, args []string) {
			flags := readFlags(cmd)
			if flags.output != "" {
				cmd.Println("Error: output flag is not supported in this command")
				return
			}

			result, err := service.Create(repo, args[0], args[1], flags.namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			cmd.Println(result)
		},
	}
}
