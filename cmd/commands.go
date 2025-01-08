package cmd

import (
	"awasm-portfolio/internal/handlers"
	"awasm-portfolio/internal/storage"
	"awasm-portfolio/internal/ui"
	"strings"

	"github.com/spf13/cobra"
)

func normalizeResourceName(resource string) string {
	plurals := map[string]string{
		"experiences": "experience",
		"skills":      "skill",
		"contacts":    "contact",
		"profiles":    "profile",
		"namespaces":  "namespace",
		"ns":          "namespace", // Add alias for namespace
	}
	if singular, exists := plurals[resource]; exists {
		return singular
	}
	return resource
}

func GetCmd(resourceManager *storage.ResourceManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [resource] [name]",
		Short: "List resources of a given type or show details of a specific resource",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			formatter := ui.TextFormatter{}
			resource := normalizeResourceName(args[0])
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			// Retrieve flags
			namespace, _ := cmd.Flags().GetString("namespace")
			allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

			// Call the handler
			output := handlers.Get(resourceManager, resource, name, namespace, allNamespaces, formatter)
			cmd.Println(output)
		},
	}

	// Set default namespace in the flag definition
	cmd.Flags().StringP("namespace", "n", "default", "Namespace for namespaced resources (default: 'default')")
	cmd.Flags().BoolP("all-namespaces", "A", false, "List resources across all namespaces")
	return cmd
}

func DescribeCmd(resourceManager *storage.ResourceManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe [resource] [name]",
		Short: "Show details of a specific resource",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			formatter := ui.TextFormatter{}
			resource := normalizeResourceName(args[0])
			name := args[1]

			// Retrieve namespace flag
			namespace, _ := cmd.Flags().GetString("namespace")

			// Call the handler
			output := handlers.Describe(resourceManager, resource, name, namespace, formatter)
			cmd.Println(output)
		},
	}

	// Set default namespace in the flag definition
	cmd.Flags().StringP("namespace", "n", "default", "Namespace for namespaced resources (default: 'default')")
	return cmd
}

func CreateCmd(resourceManager *storage.ResourceManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [resource] [name] [fields...]",
		Short: "Create a new resource",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			formatter := ui.TextFormatter{}
			resource := normalizeResourceName(args[0])
			name := args[1]
			fields := args[2:]

			// Parse fields into a map
			data := make(map[string]interface{})
			namespace, _ := cmd.Flags().GetString("namespace")
			for _, field := range fields {
				kv := strings.SplitN(field, "=", 2)
				if len(kv) == 2 {
					data[kv[0]] = kv[1]
				}
			}

			// Call the handler
			output := handlers.Create(resourceManager, resource, name, namespace, data, formatter)
			cmd.Println(output)
		},
	}

	// Set default namespace in the flag definition
	cmd.Flags().StringP("namespace", "n", "default", "Namespace for namespaced resources (default: 'default')")
	return cmd
}

func DeleteCmd(resourceManager *storage.ResourceManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [resource] [name]",
		Short: "Delete a specific resource",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			formatter := ui.TextFormatter{}
			resource := normalizeResourceName(args[0])
			name := args[1]

			// Retrieve namespace flag
			namespace, _ := cmd.Flags().GetString("namespace")

			// Call the handler
			output := handlers.Delete(resourceManager, resource, name, namespace, formatter)
			cmd.Println(output)
		},
	}

	// Set default namespace in the flag definition
	cmd.Flags().StringP("namespace", "n", "default", "Namespace for namespaced resources (default: 'default')")
	return cmd
}
