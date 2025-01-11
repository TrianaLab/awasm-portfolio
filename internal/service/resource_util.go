package service

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
	supported := SupportedResources()
	normalized := supported[strings.ToLower(resource)]
	logger.Trace(logrus.Fields{
		"input":      resource,
		"normalized": normalized,
	}, "NormalizeResourceName called")
	return normalized
}

// IsValidResource checks if a given resource kind is supported.
func IsValidResource(resource string) bool {
	_, exists := SupportedResources()[strings.ToLower(resource)]
	logger.Trace(logrus.Fields{
		"resource": resource,
		"valid":    exists,
	}, "IsValidResource called")
	return exists
}
