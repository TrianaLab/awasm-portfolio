package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"
)

type TextFormatter struct{}

func (f TextFormatter) FormatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	// Extract column headers based on the Resource interface
	headers := []string{"NAME", "NAMESPACE", "KIND"}

	var rows [][]string
	for _, resource := range resources {
		rows = append(rows, []string{
			resource.GetName(),
			resource.GetNamespace(),
			resource.GetKind(),
		})
	}

	// Determine column widths
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

	// Build the table
	var sb strings.Builder

	// Print headers
	for i, header := range headers {
		sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], header))
		if i < len(headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Print separator
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width))
		if i < len(headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], cell))
			if i < len(headers)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (f TextFormatter) FormatDetails(resource models.Resource) string {
	// Extract all fields and values using reflection
	v := reflect.ValueOf(resource).Elem()
	t := v.Type()

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Kind: %s\n", resource.GetKind()))
	sb.WriteString(fmt.Sprintf("Name: %s\n", resource.GetName()))
	sb.WriteString(fmt.Sprintf("Namespace: %s\n", resource.GetNamespace()))

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		fieldName := field.Name
		fieldValue := value.Interface()

		// Format nested structs and slices as YAML-like
		switch value.Kind() {
		case reflect.Slice:
			sb.WriteString(fmt.Sprintf("%s:\n", fieldName))
			for j := 0; j < value.Len(); j++ {
				sb.WriteString(fmt.Sprintf("  - %v\n", value.Index(j).Interface()))
			}
		case reflect.Struct:
			sb.WriteString(fmt.Sprintf("%s:\n", fieldName))
			nestedValue := reflect.ValueOf(fieldValue)
			for k := 0; k < nestedValue.NumField(); k++ {
				sb.WriteString(fmt.Sprintf("  %s: %v\n", nestedValue.Type().Field(k).Name, nestedValue.Field(k).Interface()))
			}
		default:
			sb.WriteString(fmt.Sprintf("%s: %v\n", fieldName, fieldValue))
		}
	}

	return sb.String()
}
