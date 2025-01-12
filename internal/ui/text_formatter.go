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
	for kind, group := range grouped {
		// Check if the group is for "namespace" kind
		if kind == "namespace" {
			sb.WriteString(f.formatNamespaceTable(group)) // Call the special formatting method
		} else {
			sb.WriteString(f.formatTable(group)) // Default formatting for other kinds
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (f TextFormatter) formatNamespaceTable(resources []models.Resource) string {
	if len(resources) == 0 {
		return "No resources found."
	}

	// Start with a header for the "NAME" column only
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s\n", "NAME"))

	// Print only the "Name" column for each namespace
	for _, resource := range resources {
		sb.WriteString(fmt.Sprintf("%-30s\n", resource.GetName()))
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
	headers := []string{"NAME", "NAMESPACE"} // Only add these once
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// Skip arrays and OwnerRef
		if fieldValue.Kind() == reflect.Slice || strings.ToLower(field.Name) == "ownerref" {
			continue
		}

		// Skip Name and Namespace fields since they're already included
		if strings.ToLower(field.Name) == "name" || strings.ToLower(field.Name) == "namespace" {
			continue
		}

		// Add the other fields to headers, ensuring the name is uppercase
		headers = append(headers, strings.ToUpper(field.Name))
	}

	return headers
}

func (f TextFormatter) extractRow(resource models.Resource, headers []string) []string {
	// If the resource is of kind "namespace", only include the "Name" column
	if resource.GetKind() == "namespace" {
		return []string{resource.GetName()}
	}

	row := make([]string, len(headers))

	// Extract the Name and Namespace at the beginning and skip them in the loop
	row[0] = resource.GetName()
	row[1] = resource.GetNamespace()

	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 2; i < len(headers); i++ { // Start from index 2 to skip Name and Namespace
		header := headers[i]
		found := false

		for j := 0; j < value.NumField(); j++ {
			field := typ.Field(j)
			fieldValue := value.Field(j)

			// Check if the current field matches the header
			if strings.ToUpper(field.Name) == header {
				found = true

				// Check if the field is a nested resource (implements Resource interface)
				if fieldValue.Kind() == reflect.Struct && fieldValue.Addr().Type().Implements(reflect.TypeOf((*models.Resource)(nil)).Elem()) {
					var nestedResourcePtr models.Resource
					if fieldValue.Kind() == reflect.Ptr {
						nestedResourcePtr = fieldValue.Interface().(models.Resource)
					} else {
						nestedResourcePtr = fieldValue.Addr().Interface().(models.Resource)
					}

					// Extract the name of the nested resource
					row[i] = nestedResourcePtr.GetName()
				} else if !fieldValue.IsZero() { // Handle scalar values
					row[i] = fmt.Sprintf("%v", fieldValue.Interface())
				} else { // Handle empty fields
					row[i] = ""
				}
				break
			}
		}

		// If not found, append empty value
		if !found {
			row[i] = ""
		}
	}

	return row
}

// Capitalize the field name fully
func capitalizeTableFieldName(fieldName string) string {
	return strings.ToUpper(fieldName) // Capitalize the whole name
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
