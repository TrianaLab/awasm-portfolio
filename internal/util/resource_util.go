package util

import (
	"fmt"
	"strings"
)

// SupportedResources returns a map of supported resource kinds.
// Keys are the accepted inputs, and values are the canonical names.
func SupportedResources() map[string]string {
	return map[string]string{
		"profile":        "profile",
		"profiles":       "profile",
		"namespace":      "namespace",
		"namespaces":     "namespace",
		"education":      "education",
		"educations":     "education",
		"experience":     "experience",
		"experiences":    "experience",
		"contact":        "contact",
		"contacts":       "contact",
		"certification":  "certifications",
		"certifications": "certifications",
		"contribution":   "contributions",
		"contributions":  "contributions",
		"skill":          "skills",
		"skills":         "skills",
		"ns":             "namespace", // Alias for namespace
	}
}

// NormalizeKind normalizes a given kind to its canonical form if valid,
// or returns an error if the kind is not supported.
func NormalizeKind(kind string) (string, error) {
	kind = strings.ToLower(strings.TrimSpace(kind))
	supportedResources := SupportedResources()

	// Special case for "all"
	if kind == "all" {
		return "", nil
	}

	normalized, exists := supportedResources[kind]
	if !exists {
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}
	return normalized, nil
}
