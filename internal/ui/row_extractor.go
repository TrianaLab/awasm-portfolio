package ui

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
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

				// Special handling for OwnerRef fields
				if strings.ToLower(field.Name) == "ownerref" {
					row[i] = extractOwnerRefName(fieldValue)
				} else if !fieldValue.IsZero() {
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

func extractOwnerRefName(fieldValue reflect.Value) string {
	// Check if the field is a slice
	if fieldValue.Kind() == reflect.Slice {
		if fieldValue.Len() > 0 {
			ownerRef := fieldValue.Index(0).Interface()
			if ref, ok := ownerRef.(*models.OwnerReference); ok {
				if ref.Owner != nil {
					return fmt.Sprintf("%s/%s", ref.Owner.GetKind(), ref.Owner.GetName()) // Return "kind/resource"
				}
				return fmt.Sprintf("%s/%s", ref.Kind, ref.Name) // Fallback to "kind/resource"
			}
		}
		return ""
	}

	// Handle case where OwnerRef is a single struct
	if fieldValue.Kind() == reflect.Struct {
		ownerRef := fieldValue.Interface()
		if ref, ok := ownerRef.(models.OwnerReference); ok {
			if ref.Owner != nil {
				return fmt.Sprintf("%s/%s", ref.Owner.GetKind(), ref.Owner.GetName()) // Return "kind/resource"
			}
			return fmt.Sprintf("%s/%s", ref.Kind, ref.Name) // Fallback to "kind/resource"
		}
	}

	// Log an error if the field is neither a slice nor a struct
	logger.Error(logrus.Fields{
		"fieldValueKind": fieldValue.Kind(),
	}, "Unexpected OwnerRef kind")
	return ""
}
