package main

import (
	"awasm-portfolio/cmd"
	"awasm-portfolio/internal/preload"
	"awasm-portfolio/internal/repository"
	"bytes"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	// Set trace level
	var log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.TraceLevel

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

	// Remove the root command if it's redundantly included
	if len(args) > 0 && (args[0] == "kubectl" || args[0] == "k") {
		args = args[1:]
	}

	resetFlagValues(rootCmd)

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
	})
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
	for _, subCmd := range cmd.Commands() {
		resetFlagValues(subCmd)
	}
}
