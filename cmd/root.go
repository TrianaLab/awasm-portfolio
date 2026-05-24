package cmd

import (
	"fmt"
	"sort"

	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/util"

	"github.com/spf13/cobra"
)

// NewRootCommand wires the kubectl-style root command together with the
// in-memory repository it operates on. Subcommands close over the repo
// directly — no service abstraction layer in between.
func NewRootCommand(repo *repository.InMemoryRepository) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "kubectl",
		Short:         "A CLI tool for managing resources",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.SetHelpCommand(nil)
	rootCmd.SetHelpFunc(rootHelp)

	rootCmd.AddCommand(NewCreateCommand(repo))
	rootCmd.AddCommand(NewDeleteCommand(repo))
	rootCmd.AddCommand(NewGetCommand(repo))
	rootCmd.AddCommand(NewDescribeCommand(repo))
	rootCmd.AddCommand(NewVersionCommand())

	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Specify the namespace")
	rootCmd.PersistentFlags().BoolP("all-namespaces", "A", false, "Include resources from all namespaces")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Specify the output format, either 'json' or 'yaml'")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd
}

func rootHelp(cmd *cobra.Command, _ []string) {
	cmd.Println("A CLI tool for managing resources")
	cmd.Println("\nUsage:")
	cmd.Println("  kubectl [command] [resource]")
	cmd.Println("\nAliases:")
	cmd.Println("  kubectl, k")
	cmd.Println("\nAvailable Commands:")
	cmd.Println("  create      Create a new resource")
	cmd.Println("  delete      Delete a resource")
	cmd.Println("  describe    Describe a specific resource")
	cmd.Println("  get         Get resources of a specific kind or a specific resource")
	cmd.Println("  version     Get the current version of the application")

	cmd.Println("\nAvailable Resources:")
	resources := util.SupportedResources()
	canonical := map[string]string{}
	aliases := map[string][]string{}
	for alias, kind := range resources {
		if _, exists := canonical[kind]; !exists {
			canonical[kind] = alias
		} else {
			aliases[kind] = append(aliases[kind], alias)
		}
	}
	for kind, aliasList := range aliases {
		cmd.Printf("  %s  (aliases: %s)\n", kind, fmt.Sprintf("%s", aliasList))
	}

	cmd.Println("\nFlags:")
	cmd.Println("  -A, --all-namespaces     Include resources from all namespaces")
	cmd.Println("  -h, --help               help for kubectl")
	cmd.Println("  -n, --namespace string   Specify the namespace (default \"default\")")
	cmd.Println("\nUse \"kubectl [command] --help\" for more information about a command.")
}

// canonicalResourceKinds returns the de-duplicated, sorted list of
// resource kinds, plus the "all" pseudo-kind, suitable for shell tab
// completion. Aliases are intentionally excluded — they would otherwise
// triple the candidate list and add noise.
func canonicalResourceKinds() []string {
	seen := map[string]struct{}{}
	for _, canonical := range util.SupportedResources() {
		seen[canonical] = struct{}{}
	}
	out := make([]string, 0, len(seen)+1)
	for k := range seen {
		out = append(out, k)
	}
	out = append(out, "all")
	sort.Strings(out)
	return out
}

// resourceNamesFor returns the deduplicated names of resources of the
// given kind in the active namespace context, used for completing the
// second positional argument of get/describe/delete. Dedup matters
// when the cmd context is -A (all namespaces): the same resource name
// can exist in multiple namespaces.
func resourceNamesFor(repo *repository.InMemoryRepository, kind string, namespace string) []string {
	resources, err := repo.List(kind, "", namespace)
	if err != nil {
		return nil
	}
	seen := map[string]struct{}{}
	for _, r := range resources {
		seen[r.GetName()] = struct{}{}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// completeResourceArgs is the ValidArgsFunction shared by
// get/describe/delete: the first positional completes resource kinds,
// the second completes existing resource names of that kind.
func completeResourceArgs(repo *repository.InMemoryRepository) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return canonicalResourceKinds(), cobra.ShellCompDirectiveNoFileComp
		case 1:
			flags := readFlags(cmd)
			return resourceNamesFor(repo, args[0], flags.namespace), cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}
}

// flagContext collects the persistent flag values commands care about.
type flagContext struct {
	namespace string
	output    string
}

// readFlags resolves the namespace + output flags, applying the -A
// (all-namespaces) override that clears the namespace.
func readFlags(cmd *cobra.Command) flagContext {
	namespace, _ := cmd.Flags().GetString("namespace")
	allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
	output, _ := cmd.Flags().GetString("output")
	if allNamespaces {
		namespace = ""
	}
	return flagContext{namespace: namespace, output: output}
}
