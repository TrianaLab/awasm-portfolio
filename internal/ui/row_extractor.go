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

	logger.Trace(logrus.Fields{
		"resourceName":      resource.GetName(),
		"resourceNamespace": resource.GetNamespace(),
	}, "Extracting row for resource")

	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	for i := 2; i < len(headers); i++ {
		header := headers[i]
		found := false

		for j := 0; j < value.NumField(); j++ {
			field := typ.Field(j)
			fieldValue := value.Field(j)

			logger.Trace(logrus.Fields{
				"fieldName": field.Name,
				"fieldType": fieldValue.Type().String(),
				"fieldKind": fieldValue.Kind().String(),
			}, "Inspecting field")

			if strings.ToUpper(field.Name) == header {
				found = true

				// Special handling for OwnerRef
				if strings.ToLower(field.Name) == "ownerref" {
					row[i] = formatOwnerRef(fieldValue)
					logger.Info(logrus.Fields{
						"ownerRefFormatted": row[i],
					}, "Formatted owner reference")
					break
				}

				// Handle fields of type Resource
				if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
					if fieldValue.Type().Implements(reflect.TypeOf((*models.Resource)(nil)).Elem()) {
						res := fieldValue.Interface().(models.Resource)
						row[i] = res.GetName() // Extract only the resource name
						logger.Info(logrus.Fields{
							"resourceFieldName": res.GetName(),
						}, "Resolved resource field to name")
						break
					}
				}

				// Handle nested structs (e.g., Contact, Contributions)
				if fieldValue.Kind() == reflect.Struct {
					if nestedResource, ok := fieldValue.Addr().Interface().(models.Resource); ok {
						row[i] = nestedResource.GetName() // Extract name from the nested resource
						logger.Info(logrus.Fields{
							"nestedResourceFieldName": nestedResource.GetName(),
						}, "Resolved nested struct to resource name")
						break
					}
				}

				// Handle other field types
				if fieldValue.IsValid() {
					row[i] = fmt.Sprintf("%v", fieldValue.Interface())
					logger.Trace(logrus.Fields{
						"fieldValue": fmt.Sprintf("%v", fieldValue.Interface()),
					}, "Field value set in row")
				} else {
					logger.Warn(logrus.Fields{
						"fieldName": field.Name,
					}, "Field is zero or nil")
					row[i] = ""
				}
				break
			}
		}

		if !found {
			logger.Warn(logrus.Fields{
				"header": header,
			}, "No matching field found for header")
			row[i] = ""
		}
	}

	logger.Info(logrus.Fields{
		"row": row,
	}, "Extracted row for resource")
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
