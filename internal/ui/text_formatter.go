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
	headers := f.extractHeaders(resources[0])
	var rows [][]string

	for _, resource := range resources {
		rows = append(rows, f.extractRow(resource, headers))
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

// Extract headers from the first resource, excluding OwnerRef and arrays
func (f TextFormatter) extractHeaders(resource models.Resource) []string {
	headers := []string{"NAME", "NAMESPACE"}
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// Skip arrays and OwnerRef
		if fieldValue.Kind() == reflect.Slice || strings.ToLower(field.Name) == "ownerref" {
			continue
		}

		// Add other fields to headers
		headers = append(headers, capitalizeTableFieldName(field.Name))
	}

	return headers
}

func (f TextFormatter) extractRow(resource models.Resource, headers []string) []string {
	row := []string{resource.GetName(), resource.GetNamespace()}
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 2; i < len(headers); i++ {
		header := headers[i]
		found := false

		for j := 0; j < value.NumField(); j++ {
			field := typ.Field(j)
			fieldValue := value.Field(j)

			fmt.Printf("Looking for header: %s\n", header)
			fmt.Printf("Checking field: %s, Type: %s\n", field.Name, fieldValue.Kind())

			if capitalizeTableFieldName(field.Name) == header {
				found = true

				// Check if the field is a resource (implements Resource interface)
				if fieldValue.Kind() == reflect.Struct && fieldValue.Addr().Type().Implements(reflect.TypeOf((*models.Resource)(nil)).Elem()) {
					// Check the kind of field
					fmt.Printf("Field: %s is a resource type\n", field.Name)

					// Dereference based on whether it's a pointer or not
					var nestedResourcePtr models.Resource
					if fieldValue.Kind() == reflect.Ptr {
						nestedResourcePtr = fieldValue.Interface().(models.Resource)
						fmt.Printf("Dereferenced pointer: %v\n", nestedResourcePtr)
					} else {
						nestedResourcePtr = fieldValue.Addr().Interface().(models.Resource)
						fmt.Printf("Dereferenced struct: %v\n", nestedResourcePtr)
					}

					// Extract the name of the nested resource
					name := nestedResourcePtr.GetName()
					fmt.Printf("Extracted nested resource name: %s for field: %s\n", name, header)
					row = append(row, name)
				} else if !fieldValue.IsZero() { // Scalar fields
					fmt.Printf("Extracted scalar value: %v for field: %s\n", fieldValue.Interface(), header)
					row = append(row, fmt.Sprintf("%v", fieldValue.Interface()))
				} else { // Empty field
					fmt.Printf("Field: %s is empty\n", header)
					row = append(row, "")
				}
				break
			}
		}

		// If not found in the resource, append empty value
		if !found {
			fmt.Printf("Header: %s not found in the resource\n", header)
			row = append(row, "")
		}
	}

	// Log the final row constructed for this resource
	fmt.Printf("Final row: %v\n", row)

	return row
}

// Capitalize the first letter of a field name
func capitalizeTableFieldName(fieldName string) string {
	return strings.ToUpper(fieldName[:1]) + fieldName[1:]
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
			// Log field details for debugging
			fmt.Printf("Processing field: %s\n", fieldName)
			fmt.Printf("Field value type: %s, Kind: %s\n", fieldValue.Type(), fieldValue.Kind())

			// Dereference the struct pointer properly
			var nestedResourcePtr models.Resource
			if fieldValue.Kind() == reflect.Ptr {
				nestedResourcePtr = fieldValue.Interface().(models.Resource)
				fmt.Printf("Dereferenced pointer: %v\n", nestedResourcePtr)
			} else {
				nestedResourcePtr = fieldValue.Addr().Interface().(models.Resource)
				fmt.Printf("Dereferenced struct: %v\n", nestedResourcePtr)
			}

			// Check if the Name field is properly populated
			name := nestedResourcePtr.GetName()
			fmt.Printf("Extracted Name using GetName(): %s\n", name)

			if name == "" {
				fmt.Printf("Warning: Name is empty for field: %s\n", fieldName)
			}

			buffer.WriteString(fmt.Sprintf("%s%s: %s\n", indentation(indent), fieldName, name))
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
