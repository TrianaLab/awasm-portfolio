package service

import (
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

// NormalizeResourceName normalizes a given resource name to its canonical form.
func NormalizeResourceName(resource string) string {
	supported := SupportedResources()
	if canonical, exists := supported[strings.ToLower(resource)]; exists {
		return canonical
	}
	return resource
}

// IsValidResource checks if a given resource kind is supported.
func IsValidResource(resource string) bool {
	_, exists := SupportedResources()[strings.ToLower(resource)]
	return exists
}
