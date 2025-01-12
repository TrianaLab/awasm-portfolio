package ui

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"strings"

	"github.com/sirupsen/logrus"
)

// TableFormatter formats resources into tables
type TableFormatter struct{}

// FormatTable formats resources into grouped tables
func (f TableFormatter) FormatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	grouped := groupResourcesByKind(resources)

	var sb strings.Builder
	for kind, group := range grouped {
		if kind == "namespace" {
			sb.WriteString(f.formatNamespaceTable(group))
		} else {
			sb.WriteString(f.formatGenericTable(group))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatNamespaceTable formats namespace resources as a single-column table
// formatNamespaceTable formats namespace resources with AGE column
func (f TableFormatter) formatNamespaceTable(resources []models.Resource) string {
	var sb strings.Builder
	headers := []string{"NAME", "AGE"} // Add AGE column
	rows := [][]string{}

	for _, resource := range resources {
		name := resource.GetName()
		age := calculateAge(resource.GetCreationTimestamp())
		rows = append(rows, []string{name, age})
		logger.Trace(logrus.Fields{
			"name": name,
			"age":  age,
		}, "Processed namespace resource")
	}

	colWidths := calculateColumnWidths(headers, rows)
	formatHeaders(&sb, headers, colWidths)
	formatRows(&sb, rows, colWidths)

	return sb.String()
}

// formatGenericTable formats a generic table for resources
func (f TableFormatter) formatGenericTable(resources []models.Resource) string {
	headers := extractHeaders(resources[0])
	rows := extractRows(resources, headers)
	colWidths := calculateColumnWidths(headers, rows)

	var sb strings.Builder
	formatHeaders(&sb, headers, colWidths)
	formatRows(&sb, rows, colWidths)
	return sb.String()
}
