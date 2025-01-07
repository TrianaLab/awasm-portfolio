package main

import (
	"awasm-portfolio/cmd"
	"awasm-portfolio/internal/storage"
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	// Initialize the resource manager
	resourceManager := storage.NewResourceManager()

	// Preload data
	storage.PreloadData(resourceManager)

	// Initialize commands
	rootCmd := cmd.RootCmd()
	rootCmd.AddCommand(cmd.GetCmd(resourceManager))
	rootCmd.AddCommand(cmd.DescribeCmd(resourceManager))
	rootCmd.AddCommand(cmd.CreateCmd(resourceManager))
	rootCmd.AddCommand(cmd.DeleteCmd(resourceManager))

	// Expose executeCommand to JavaScript
	js.Global().Set("executeCommand", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return "Error: No command provided"
		}
		command := args[0].String()
		return runCLICommand(rootCmd, command)
	}))

	// Keep Go running
	select {}
}

func runCLICommand(rootCmd *cobra.Command, command string) string {
	args := strings.Fields(command)

	// Strip "kubectl" or "k" prefix if present
	if len(args) > 0 && (args[0] == "kubectl" || args[0] == "k") {
		args = args[1:]
	}

	// Reset flag values recursively
	resetFlagValues(rootCmd)

	// Set the command arguments
	rootCmd.SetArgs(args)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	if err := rootCmd.Execute(); err != nil {
		buf.WriteString(fmt.Sprintf("Error: %s", err.Error()))
	}

	return strings.TrimSpace(buf.String())
}

// resetFlagValues resets all flag values to their defaults for the root command and its subcommands.
func resetFlagValues(cmd *cobra.Command) {
	// Reset flags for the current command
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})

	// Recursively reset flags for subcommands
	for _, subCmd := range cmd.Commands() {
		resetFlagValues(subCmd)
	}
}
