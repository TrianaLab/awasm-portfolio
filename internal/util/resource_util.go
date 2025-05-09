package util

import (
	"fmt"
	"strings"
)

// SupportedResources returns a map of supported resource kinds.
// Keys are the accepted inputs, and values are the canonical names.
func SupportedResources() map[string]string {
	return map[string]string{
		"resume":       "resume",
		"resumes":      "resume",
		"basic":        "basics",
		"basics":       "basics",
		"namespace":    "namespace",
		"namespaces":   "namespace",
		"ns":           "namespace",
		"work":         "work",
		"works":        "work",
		"volunteer":    "volunteer",
		"volunteers":   "volunteer",
		"education":    "education",
		"educations":   "education",
		"award":        "award",
		"awards":       "award",
		"certificate":  "certificate",
		"certificates": "certificate",
		"publication":  "publication",
		"publications": "publication",
		"skill":        "skill",
		"skills":       "skill",
		"language":     "language",
		"languages":    "language",
		"interest":     "interest",
		"interests":    "interest",
		"reference":    "reference",
		"references":   "reference",
		"project":      "project",
		"projects":     "project",
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
