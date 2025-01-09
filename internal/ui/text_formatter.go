package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"strings"
)

type TextFormatter struct{}

func (f TextFormatter) FormatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found.\n"
	}

	var rows []string
	for _, res := range resources {
		rows = append(rows, fmt.Sprintf("%s\t%s", res.GetName(), res.GetNamespace()))
	}

	return strings.Join(rows, "\n")
}

func (f TextFormatter) FormatDetails(resource models.Resource) string {
	details := fmt.Sprintf(
		"Name: %s\nNamespace: %s\n",
		resource.GetName(),
		resource.GetNamespace(),
	)

	// Optionally format owner references
	if owners := resource.GetOwnerReferences(); len(owners) > 0 {
		details += "Owner References:\n"
		for _, owner := range owners {
			details += fmt.Sprintf("  - Kind: %s, Name: %s\n", owner.Kind, owner.Name)
		}
	}

	return details
}
