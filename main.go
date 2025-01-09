package main

import (
	"awasm-portfolio/cmd"
	"awasm-portfolio/internal/preload"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/service"
	"awasm-portfolio/internal/ui"
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	// Initialize repository
	repo := repository.NewInMemoryRepository()

	// Initialize formatter
	formatter := ui.TextFormatter{}

	// Initialize service with repository and formatter
	resourceService := service.NewResourceService(repo, formatter)

	// Preload data
	preload.PreloadData(repo)

	// Initialize commands
	rootCmd := cmd.NewRootCommand(resourceService)

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
		if strings.Contains(err.Error(), "unknown command") || strings.Contains(err.Error(), "requires at least") {
			buf.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		} else {
			buf.WriteString(err.Error())
		}
	}

	return strings.TrimSpace(buf.String())
}

func resetFlagValues(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
	for _, subCmd := range cmd.Commands() {
		resetFlagValues(subCmd)
	}
}
