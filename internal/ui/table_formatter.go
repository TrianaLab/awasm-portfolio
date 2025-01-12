package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"strings"
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
			sb.WriteString(f.formatTable(group))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatNamespaceTable formats namespace resources as a single-column table
func (f TableFormatter) formatNamespaceTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s\n", "NAME"))
	for _, resource := range resources {
		sb.WriteString(fmt.Sprintf("%-30s\n", resource.GetName()))
	}

	return sb.String()
}

// formatTable formats a generic table for resources
func (f TableFormatter) formatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	headers := extractHeaders(resources[0])
	rows := extractRows(resources, headers)
	colWidths := calculateColumnWidths(headers, rows)

	var sb strings.Builder
	formatHeaders(&sb, headers, colWidths)
	formatRows(&sb, rows, colWidths)

	return sb.String()
}
