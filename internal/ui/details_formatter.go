package ui

import (
	"awasm-portfolio/internal/models"
)

// FormatDetails formats a resource into YAML-like details
func FormatDetails(resources []models.Resource) string {
	return formatAsYAML(resources)
}
