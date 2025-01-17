package ui

import (
	"awasm-portfolio/internal/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// calculateColumnWidths calculates the width for each column
func calculateColumnWidths(headers []string, rows [][]string) []int {
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	return colWidths
}

// formatHeaders formats the table headers
func formatHeaders(sb *strings.Builder, headers []string, colWidths []int) {
	for i, header := range headers {
		sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], header))
		if i < len(headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")
}

// formatRows formats the table rows
func formatRows(sb *strings.Builder, rows [][]string, colWidths []int) {
	for _, row := range rows {
		for i, cell := range row {
			sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], cell))
			if i < len(row)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}
}

// groupResourcesByKind groups resources by their kind
func groupResourcesByKind(resources []models.Resource) map[string][]models.Resource {
	grouped := make(map[string][]models.Resource)
	for _, resource := range resources {
		kind := resource.GetKind()
		grouped[kind] = append(grouped[kind], resource)
	}
	return grouped
}

// calculateAge calculates the age of a resource
func calculateAge(timestamp time.Time) string {
	if timestamp.IsZero() {
		return ""
	}
	duration := time.Since(timestamp)
	switch {
	case duration.Hours() >= 24:
		return fmt.Sprintf("%dd", int(duration.Hours()/24))
	case duration.Hours() >= 1:
		return fmt.Sprintf("%dh", int(duration.Hours()))
	case duration.Minutes() >= 1:
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	default:
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
}

// extractRowsFromSchema generates rows based on schema extractors
func extractRowsFromSchema(resources []models.Resource, schema Schema) [][]string {
	rows := make([][]string, len(resources))
	for i, resource := range resources {
		row := make([]string, len(schema.Extractors))
		for j, extractor := range schema.Extractors {
			row[j] = extractor(resource)
		}
		rows[i] = row
	}
	return rows
}

// formatAsTable formats resources as a table
func formatAsTable(resources []models.Resource, schemas map[string]Schema) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	grouped := groupResourcesByKind(resources)
	var sb strings.Builder

	for kind, group := range grouped {
		// Get schema for the resource kind
		schema, exists := schemas[kind]
		if !exists {
			schema = schemas["default"]
		}

		headers := schema.Headers
		rows := extractRowsFromSchema(group, schema)
		colWidths := calculateColumnWidths(headers, rows)

		// Format headers and rows
		formatHeaders(&sb, headers, colWidths)
		formatRows(&sb, rows, colWidths)
		sb.WriteString("\n")
	}

	return sb.String()
}

// formatAsJSON marshals resources into JSON
func formatAsJSON(resources []models.Resource) string {
	data, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return "Error formatting resources as JSON: " + err.Error()
	}
	return string(data)
}

// formatAsYAML marshals resources into YAML
func formatAsYAML(resources []models.Resource) string {
	data, err := yaml.Marshal(resources)
	if err != nil {
		return "Error formatting resources as YAML: " + err.Error()
	}
	return string(data)
}
