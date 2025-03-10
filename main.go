//go:build js && wasm
// +build js,wasm

package main

import (
	"awasm-portfolio/cmd"
	"awasm-portfolio/internal/preload"
	"awasm-portfolio/internal/repository"
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {

	// Initialize repository and factory
	repo := repository.NewInMemoryRepository()

	// Preload data
	preload.PreloadData(repo)

	// Initialize commands
	rootCmd := cmd.NewRootCommand(repo)

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

	if len(args) == 0 || (args[0] != "kubectl" && args[0] != "k") {
		return fmt.Sprintf("Error: unknown command '%s'", args[0])
	}

	// Remove the root command prefix
	args = args[1:]

	// Reset flag values
	resetFlagValues(rootCmd)

	// Set the command arguments
	rootCmd.SetArgs(args)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	if err := rootCmd.Execute(); err != nil {
		buf.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
	}

	return strings.TrimSpace(buf.String())
}

func resetFlagValues(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})
	for _, subCmd := range cmd.Commands() {
		resetFlagValues(subCmd)
	}
}
