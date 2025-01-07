package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A CLI tool for managing resources in the portfolio",
	Long:  "This CLI tool emulates Kubernetes-like commands for managing resources such as profiles, contacts, and more.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Error: Please specify a valid subcommand")
	},
}

func RootCmd() *cobra.Command {
	return rootCmd
}

func RunCommand(command string) string {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	args := strings.Fields(command)
	if len(args) > 0 && (args[0] == "kubectl" || args[0] == "k") {
		args = args[1:]
	}
	rootCmd.SetArgs(args)

	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})

	if err := rootCmd.Execute(); err != nil {
		buf.WriteString(fmt.Sprintf("Error: %s", err.Error()))
	}

	return strings.TrimSpace(buf.String())
}
