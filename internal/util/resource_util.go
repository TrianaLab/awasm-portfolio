package util

import (
	"awasm-portfolio/internal/logger"
	"strings"

	"github.com/sirupsen/logrus"
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
	resource = strings.ToLower(resource)
	supported := SupportedResources()
	normalized, exists := supported[resource]
	if !exists {
		logger.Trace(logrus.Fields{
			"input": resource,
		}, "Unsupported resource kind in NormalizeResourceName")
		return resource // Return the original kind to preserve it in error messages
	}
	return normalized
}

// IsValidResource checks if a given resource kind is supported.
func IsValidResource(resource string) bool {
	resource = strings.ToLower(resource)
	_, exists := SupportedResources()[resource]
	return exists
}
