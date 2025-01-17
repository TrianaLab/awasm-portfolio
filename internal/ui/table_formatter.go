package ui

import (
	"awasm-portfolio/internal/models"
	"strings"
)

// TableFormatter formats resources into tables
type TableFormatter struct {
	schemas map[string]Schema
}

// NewTableFormatter creates a new TableFormatter
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{
		schemas: GenerateSchemas(),
	}
}

// FormatTable formats resources dynamically
func (f *TableFormatter) FormatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	// Group resources by kind
	grouped := groupResourcesByKind(resources)
	var sb strings.Builder

	for kind, group := range grouped {
		// Get schema for the resource kind
		schema, exists := f.schemas[kind]
		if !exists {
			schema = f.schemas["default"]
		}

		// Extract headers and rows
		headers := schema.Headers
		rows := extractRowsFromSchema(group, schema)

		// Calculate column widths and format table
		colWidths := calculateColumnWidths(headers, rows)
		formatHeaders(&sb, headers, colWidths)
		formatRows(&sb, rows, colWidths)
		sb.WriteString("\n")
	}

	return sb.String()
}
