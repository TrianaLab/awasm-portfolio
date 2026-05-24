//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/TrianaLab/awasm-portfolio/cmd"
	"github.com/TrianaLab/awasm-portfolio/internal/preload"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"

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

	// Expose completeCommand to JavaScript. Returns a JSON-encoded
	// {"candidates":[...]} so the worker can hand a plain string back to
	// the main thread (matching the executeCommand round-trip shape).
	js.Global().Set("completeCommand", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return `{"candidates":[]}`
		}
		return completeCLICommand(rootCmd, args[0].String())
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

// completeCLICommand drives Cobra's hidden __complete command and returns
// the candidate list (without descriptions or trailing directive lines)
// as a JSON object: {"candidates":["get","describe",...]}.
//
// The line is what the user has typed so far, including the binary
// prefix. A trailing space means "complete the next, empty token";
// otherwise we complete the final partial token in place.
func completeCLICommand(rootCmd *cobra.Command, line string) string {
	fields := strings.Fields(line)
	endsWithSpace := strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t")

	// No binary typed yet → offer the two entry points.
	if len(fields) == 0 {
		return encodeCandidates([]string{"kubectl"})
	}
	if fields[0] != "kubectl" && fields[0] != "k" {
		if len(fields) == 1 && !endsWithSpace {
			return encodeCandidates(filterPrefix([]string{"kubectl", "k"}, fields[0]))
		}
		return encodeCandidates(nil)
	}

	// Drop the binary; cobra's __complete sees the post-binary args.
	completionArgs := fields[1:]
	if endsWithSpace {
		completionArgs = append(completionArgs, "")
	}

	resetFlagValues(rootCmd)
	rootCmd.SetArgs(append([]string{cobra.ShellCompRequestCmd}, completionArgs...))

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	_ = rootCmd.Execute()

	return encodeCandidates(parseCobraCompletions(buf.String()))
}

// parseCobraCompletions extracts the candidate tokens from the raw
// __complete output. The format is one candidate per line (optionally
// "candidate\tdescription"), terminated by a ":<directive>" sentinel
// line followed by free-form human text we ignore.
func parseCobraCompletions(out string) []string {
	var candidates []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ":") {
			break
		}
		if tab := strings.IndexByte(line, '\t'); tab >= 0 {
			line = line[:tab]
		}
		candidates = append(candidates, line)
	}
	return candidates
}

func filterPrefix(items []string, prefix string) []string {
	var out []string
	for _, s := range items {
		if strings.HasPrefix(s, prefix) {
			out = append(out, s)
		}
	}
	return out
}

func encodeCandidates(candidates []string) string {
	if candidates == nil {
		candidates = []string{}
	}
	payload, err := json.Marshal(map[string][]string{"candidates": candidates})
	if err != nil {
		return `{"candidates":[]}`
	}
	return string(payload)
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
