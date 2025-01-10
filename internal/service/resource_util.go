package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"reflect"
	"strings"
)

// SupportedResources dynamically retrieves resource types and aliases
func SupportedResources() map[string]string {
	// Map to hold singular -> singular and alias/plural -> singular mappings
	supported := map[string]string{}

	// List of model types implementing the Resource interface
	types := []models.Resource{
		&types.Profile{},
		&types.Namespace{},
		&types.Education{},
		&types.Experience{},
		&types.Contact{},
		&types.Certifications{},
		&types.Contributions{},
		&types.Skills{},
	}

	// Populate supported map dynamically
	for _, t := range types {
		kind := normalizeKind(reflect.TypeOf(t).Elem().Name()) // Get singular form
		supported[kind] = kind                                 // Add singular form

		// Add plural form
		plural := kind + "s" // Simple pluralization rule
		supported[plural] = kind
	}

	// Add custom aliases
	supported["ns"] = "namespace"

	return supported
}

// NormalizeResourceName converts aliases/plurals to singular
func NormalizeResourceName(resource string) string {
	if strings.EqualFold(resource, "all") {
		return "all" // Special case for "all"
	}

	supported := SupportedResources()

	// Check if the resource is in the supported map
	if singular, exists := supported[resource]; exists {
		return singular
	}

	// Return the original if not recognized
	return resource
}

// IsValidResource checks if a given resource type is supported
func IsValidResource(resource string) bool {
	_, exists := SupportedResources()[resource]
	return exists
}

// normalizeKind converts struct names to lowercase singular form
func normalizeKind(kind string) string {
	return strings.ToLower(kind)
}
