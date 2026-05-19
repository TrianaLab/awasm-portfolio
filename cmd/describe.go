package cmd

import (
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/service"

	"github.com/spf13/cobra"
)

func NewDescribeCommand(repo *repository.InMemoryRepository) *cobra.Command {
	return &cobra.Command{
		Use:   "describe [kind] [name]",
		Short: "Describe a specific resource",
		Args:  cobra.RangeArgs(0, 2),
		Example: `
# Describe the profile john-doe in the dev namespace
kubectl describe profile john-doe -n dev
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}

			flags := readFlags(cmd)
			if flags.output != "" {
				cmd.Println("Error: output flag is not supported in this command")
				return
			}

			kind := args[0]
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			result, err := service.Describe(repo, kind, name, flags.namespace)
			if err != nil {
				cmd.Println("Error: ", err)
				return
			}
			cmd.Println(result)
		},
	}
}
