package ui

import (
	"awasm-portfolio/internal/models"
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// DetailsFormatter formats resources into a detailed YAML-like structure
type DetailsFormatter struct{}

// FormatDetails formats a resource into YAML-like details
func (f DetailsFormatter) FormatDetails(resource models.Resource) string {
	var buffer bytes.Buffer
	formatResource(&buffer, resource, 0, true)
	return buffer.String()
}

// formatResource recursively formats a resource and its fields
func formatResource(buffer *bytes.Buffer, resource interface{}, indent int, includeMeta bool) {
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	indentation := strings.Repeat("  ", indent)
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		if !includeMeta && (field.Name == "Name" || field.Name == "Namespace") {
			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			buffer.WriteString(fmt.Sprintf("%s%s:\n", indentation, field.Name))
			formatResource(buffer, fieldValue.Addr().Interface(), indent+1, false)
		} else if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("%s%s:\n", indentation, field.Name))
			for j := 0; j < fieldValue.Len(); j++ {
				formatResource(buffer, fieldValue.Index(j).Addr().Interface(), indent+1, true)
			}
		} else if !fieldValue.IsZero() {
			buffer.WriteString(fmt.Sprintf("%s%s: %v\n", indentation, field.Name, fieldValue.Interface()))
		}
	}
}
