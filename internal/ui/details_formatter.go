package ui

import (
	"awasm-portfolio/internal/models"

	"gopkg.in/yaml.v3"
)

// DetailsFormatter formats resources into a detailed YAML-like structure
type DetailsFormatter struct{}

// FormatDetails formats a resource into YAML-like details
func (f DetailsFormatter) FormatDetails(resource models.Resource) string {
	// Convert the resource to YAML format
	data, err := yaml.Marshal(resource)
	if err != nil {
		return "Error formatting resource: " + err.Error()
	}
	return string(data)
}
