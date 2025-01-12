package ui

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

// capitalizeFieldName capitalizes the first letter of a field name
func capitalizeFieldName(fieldName string) string {
	return strings.ToUpper(fieldName[:1]) + strings.ToLower(fieldName[1:])
}

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

func summarizeArray(fieldValue reflect.Value, header string) string {
	count := fieldValue.Len()
	if count == 0 {
		return ""
	}

	// Use the header name to create a meaningful summary
	label := strings.ToLower(strings.TrimSuffix(header, "S")) // Singular form of header
	if count > 1 {
		label += "s" // Pluralize for counts > 1
	}

	return fmt.Sprintf("%d %s", count, label)
}

func formatKindResource(fieldValue reflect.Value) string {
	// Handle pointers
	if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
		if resource, ok := fieldValue.Interface().(models.Resource); ok {
			if resource.GetKind() != "" && resource.GetName() != "" {
				return fmt.Sprintf("%s/%s", resource.GetKind(), resource.GetName())
			}
			return "" // Return empty string for incomplete resource
		}
	}

	// Handle structs
	if fieldValue.Kind() == reflect.Struct {
		// Check if the field is an OwnerReference
		if ownerRef, ok := fieldValue.Interface().(models.OwnerReference); ok {
			if ownerRef.Kind != "" && ownerRef.Name != "" {
				return fmt.Sprintf("%s/%s", ownerRef.Kind, ownerRef.Name)
			}
			return "" // Return empty string for incomplete OwnerReference
		}

		// Check if the field implements the Resource interface
		if resource, ok := fieldValue.Addr().Interface().(models.Resource); ok {
			if resource.GetKind() != "" && resource.GetName() != "" {
				return fmt.Sprintf("%s/%s", resource.GetKind(), resource.GetName())
			}
			return "" // Return empty string for incomplete resource
		}
	}

	// Log a warning if the field cannot be formatted
	logger.Warn(logrus.Fields{
		"fieldKind": fieldValue.Kind(),
		"fieldType": fieldValue.Type().String(),
	}, "Unable to format field as kind/name")
	return ""
}
