package ui

import (
	"awasm-portfolio/internal/models"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// DetailsFormatter formats resources into a detailed YAML-like structure
type DetailsFormatter struct{}

// FormatDetails formats a resource into YAML-like details
func (f DetailsFormatter) FormatDetails(resource models.Resource) string {
	var buffer bytes.Buffer
	formatResource(&buffer, resource, 0, true) // Pass `true` to indicate it's the top-level resource
	return buffer.String()
}

func formatResource(buffer *bytes.Buffer, resource interface{}, indent int, isTopLevel bool) {
	value := reflect.ValueOf(resource).Elem()
	typ := value.Type()

	indentation := strings.Repeat("  ", indent)

	// Handle top-level meta fields
	if isTopLevel {

		if field := value.FieldByName("Name"); field.IsValid() && !field.IsZero() {
			buffer.WriteString(fmt.Sprintf("%sName: %v\n", indentation, field.Interface()))
		}
		if field := value.FieldByName("Namespace"); field.IsValid() && !field.IsZero() {
			buffer.WriteString(fmt.Sprintf("%sNamespace: %v\n", indentation, field.Interface()))
		}
		if field := value.FieldByName("OwnerRef"); field.IsValid() && !field.IsZero() {
			buffer.WriteString(fmt.Sprintf("%sOwnerRef:\n", indentation))
			formatResource(buffer, field.Addr().Interface(), indent+1, false)
		}
		if field := value.FieldByName("CreationTimestamp"); field.IsValid() && !field.IsZero() {
			timestamp := field.Interface().(time.Time)
			buffer.WriteString(fmt.Sprintf("%sCreationTimestamp: %s\n", indentation, timestamp.Format("2006-01-02 15:04:05")))
		}
	}

	// Iterate over struct fields
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldValue := value.Field(i)

		// Skip top-level meta fields for the root resource
		if isTopLevel && (field.Name == "Name" || field.Name == "Namespace" || field.Name == "OwnerRef" || field.Name == "CreationTimestamp") {
			continue
		}

		// Skip meta fields for child resources
		if !isTopLevel && (field.Name == "Name" || field.Name == "Namespace" || field.Name == "OwnerRef" || field.Name == "CreationTimestamp") {
			continue
		}

		if fieldValue.Kind() == reflect.Struct && !fieldValue.IsZero() {
			buffer.WriteString(fmt.Sprintf("%s%s:\n", indentation, field.Name))
			formatResource(buffer, fieldValue.Addr().Interface(), indent+1, false)
		} else if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
			// Write the field header once for slices
			buffer.WriteString(fmt.Sprintf("%s%s:\n", indentation, field.Name))

			for j := 0; j < fieldValue.Len(); j++ {
				// Pass the child element directly with increased indentation
				formatResource(buffer, fieldValue.Index(j).Addr().Interface(), indent+1, false)
			}
		} else if !fieldValue.IsZero() {
			buffer.WriteString(fmt.Sprintf("%s%s: %v\n", indentation, field.Name, fieldValue.Interface()))
		}
	}
}
