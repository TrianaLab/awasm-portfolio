package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"
	"time"
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

// summarizeArray returns a summary of the array or slice
func summarizeArray(fieldValue reflect.Value) string {
	length := fieldValue.Len()
	if length > 0 {
		return fmt.Sprintf("%d items", length)
	}
	return "0 items"
}

// formatNestedField formats a nested field as "kind/name" if applicable
func formatNestedField(fieldValue reflect.Value) string {
	if !fieldValue.IsValid() || (fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil()) {
		return ""
	}

	if ownerRef, ok := fieldValue.Interface().(models.OwnerReference); ok {
		if ownerRef.Name != "" {
			return fmt.Sprintf("%s/%s", ownerRef.Kind, ownerRef.Name)
		}
	}

	if resource, ok := fieldValue.Interface().(models.Resource); ok {
		formatted := fmt.Sprintf("%s/%s", resource.GetKind(), resource.GetName())
		return formatted
	}

	if fieldValue.Kind() == reflect.Struct {
		if nestedField := fieldValue.Addr(); nestedField.IsValid() {
			if resource, ok := nestedField.Interface().(models.Resource); ok {
				if resource.GetName() == "" {
					return ""
				}
				formatted := fmt.Sprintf("%s/%s", resource.GetKind(), resource.GetName())
				return formatted
			}
		}
	}

	return fmt.Sprintf("%v", fieldValue.Interface())
}
