package ui

import (
	"awasm-portfolio/internal/models"
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type TextFormatter struct{}

func (f TextFormatter) FormatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	// Group resources by kind
	grouped := make(map[string][]models.Resource)
	for _, resource := range resources {
		grouped[resource.GetKind()] = append(grouped[resource.GetKind()], resource)
	}

	var sb strings.Builder

	// Iterate over each group
	for _, group := range grouped {
		sb.WriteString(f.formatTable(group))
		sb.WriteString("\n")
	}

	return sb.String()
}

// Helper function to format a single table of resources
func (f TextFormatter) formatTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	// Extract column headers
	headers := []string{"NAME", "NAMESPACE"}
	var rows [][]string

	for _, resource := range resources {
		row := []string{resource.GetName(), resource.GetNamespace()}
		rows = append(rows, row)
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

// FormatDetails formats a resource into a clean YAML-like structure.
func FormatDetails(resource interface{}) string {
	var buffer bytes.Buffer
	formatResource(&buffer, resource, 0, true)
	return buffer.String()
}

// formatResource handles formatting for a single resource, including nested resources and arrays.
func formatResource(buffer *bytes.Buffer, resource interface{}, indent int, includeMeta bool) {
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	indentation := func(level int) string {
		return strings.Repeat("  ", level)
	}

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)
		fieldName := capitalizeFieldName(field.Name)

		// Skip metadata fields for nested resources
		if !includeMeta && (strings.ToLower(fieldName) == "name" || strings.ToLower(fieldName) == "namespace" || strings.ToLower(fieldName) == "ownerref") {
			continue
		}

		// Check if the field is a resource (implements Resource interface)
		if fieldValue.Kind() == reflect.Struct && fieldValue.Addr().Type().Implements(reflect.TypeOf((*models.Resource)(nil)).Elem()) {
			resourceName := fieldValue.FieldByName("Name").String()
			buffer.WriteString(fmt.Sprintf("%s%s: %s\n", indentation(indent), fieldName, resourceName))
			formatResource(buffer, fieldValue.Addr().Interface(), indent+1, false)
			continue
		}

		// Handle slices
		if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
			for j := 0; j < fieldValue.Len(); j++ {
				element := fieldValue.Index(j).Interface()
				buffer.WriteString(fmt.Sprintf("%s- ", indentation(indent)))
				if reflect.ValueOf(element).Kind() == reflect.Struct {
					formatResource(buffer, fieldValue.Index(j).Addr().Interface(), indent+1, true)
				} else {
					buffer.WriteString(fmt.Sprintf("%v\n", element))
				}
			}
			continue
		}

		// Print scalar fields
		if !fieldValue.IsZero() {
			buffer.WriteString(fmt.Sprintf("%s%s: %v\n", indentation(indent), fieldName, fieldValue.Interface()))
		}
	}
}

// capitalizeFieldName capitalizes the first letter of a field name.
func capitalizeFieldName(fieldName string) string {
	return strings.ToUpper(fieldName[:1]) + strings.ToLower(fieldName[1:])
}
