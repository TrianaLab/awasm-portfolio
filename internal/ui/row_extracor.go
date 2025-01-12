package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"
)

// extractHeaders extracts column headers from a resource
func extractHeaders(resource models.Resource) []string {
	headers := []string{"NAME", "NAMESPACE"}
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		if field.Name != "Name" && field.Name != "Namespace" {
			headers = append(headers, strings.ToUpper(field.Name))
		}
	}

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

// extractRow extracts a single row for a resource
func extractRow(resource models.Resource, headers []string) []string {
	row := make([]string, len(headers))
	row[0] = resource.GetName()
	row[1] = resource.GetNamespace()

	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 2; i < len(headers); i++ {
		header := headers[i]
		for j := 0; j < value.NumField(); j++ {
			field := typ.Field(j)
			if strings.ToUpper(field.Name) == header {
				fieldValue := value.Field(j)
				if !fieldValue.IsZero() {
					row[i] = fmt.Sprintf("%v", fieldValue.Interface())
				} else {
					row[i] = ""
				}
				break
			}
		}
	}

	return row
}
