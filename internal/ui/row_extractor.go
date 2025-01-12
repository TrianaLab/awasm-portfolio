package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"
)

// extractHeaders extracts column headers from a resource
func extractHeaders(resource models.Resource) []string {
	headers := []string{"NAME", "NAMESPACE", "OWNERREF"} // Base headers
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		if field.Name != "Name" && field.Name != "Namespace" && field.Name != "OwnerRef" && field.Name != "Age" && field.Name != "CreationTimestamp" {
			headers = append(headers, strings.ToUpper(field.Name))
		}
	}

	headers = append(headers, "AGE") // Ensure AGE is always last
	return headers
}

// extractRows extracts rows for a table based on headers
func extractRows(resources []models.Resource, headers []string) [][]string {
	rows := [][]string{}
	for _, resource := range resources {
		row := extractRow(resource, headers)
		rows = append(rows, row)
	}
	return rows
}

// extractRow extracts a single row from a resource
func extractRow(resource models.Resource, headers []string) []string {
	row := make([]string, len(headers))

	for i, header := range headers {
		switch strings.ToLower(header) {
		case "age":
			age := calculateAge(resource.GetCreationTimestamp())
			row[i] = age
		case "creationtimestamp":
			row[i] = "" // Ensure CREATIONTIMESTAMP is always empty
		default:
			fieldName := strings.Title(strings.ToLower(header))
			fieldValue := reflect.ValueOf(resource).Elem().FieldByName(fieldName)

			if fieldValue.IsValid() {
				switch fieldValue.Kind() {
				case reflect.Slice:
					row[i] = summarizeArray(fieldValue) // Summarize arrays
				case reflect.Ptr, reflect.Struct:
					row[i] = formatNestedField(fieldValue)
				default:
					row[i] = fmt.Sprintf("%v", fieldValue.Interface())
				}
			} else {
				row[i] = ""
			}
		}
	}

	return row
}
