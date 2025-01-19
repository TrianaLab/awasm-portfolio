package ui

import (
	"awasm-portfolio/internal/models"
	"strings"
)

// FormatTable formats resources dynamically
func FormatTable(resources []models.Resource, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return formatAsJSON(resources)
	case "yaml":
		return formatAsYAML(resources)
	default:
		return formatAsTable(resources, GenerateSchemas())
	}
}
