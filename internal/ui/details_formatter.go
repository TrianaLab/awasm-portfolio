package ui

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

	logger.Trace(logrus.Fields{
		"resource":   typ.Name(),
		"isTopLevel": isTopLevel,
	}, "Starting formatResource")

	// Handle top-level meta fields
	if isTopLevel {
		logger.Trace(logrus.Fields{
			"resource": typ.Name(),
		}, "Processing top-level meta fields")

		if field := value.FieldByName("Name"); field.IsValid() && !field.IsZero() {
			buffer.WriteString(fmt.Sprintf("%sName: %v\n", indentation, field.Interface()))
			logger.Trace(logrus.Fields{"field": "Name", "value": field.Interface()}, "Processed Name")
		}
		if field := value.FieldByName("Namespace"); field.IsValid() && !field.IsZero() {
			buffer.WriteString(fmt.Sprintf("%sNamespace: %v\n", indentation, field.Interface()))
			logger.Trace(logrus.Fields{"field": "Namespace", "value": field.Interface()}, "Processed Namespace")
		}
		if field := value.FieldByName("OwnerRef"); field.IsValid() && !field.IsZero() {
			buffer.WriteString(fmt.Sprintf("%sOwnerRef:\n", indentation))
			formatResource(buffer, field.Addr().Interface(), indent+1, false)
			logger.Trace(logrus.Fields{"field": "OwnerRef", "value": field.Interface()}, "Processed OwnerRef")
		}
		if field := value.FieldByName("CreationTimestamp"); field.IsValid() && !field.IsZero() {
			timestamp := field.Interface().(time.Time)
			buffer.WriteString(fmt.Sprintf("%sCreationTimestamp: %s\n", indentation, timestamp.Format("2006-01-02 15:04:05")))
			logger.Trace(logrus.Fields{"field": "CreationTimestamp", "value": timestamp}, "Processed CreationTimestamp")
		}
	}

	// Iterate over struct fields
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			logger.Warn(logrus.Fields{"field": field.Name}, "Skipping unexported field")
			continue
		}

		fieldValue := value.Field(i)

		// Skip top-level meta fields for the root resource
		if isTopLevel && (field.Name == "Name" || field.Name == "Namespace" || field.Name == "OwnerRef" || field.Name == "CreationTimestamp") {
			logger.Trace(logrus.Fields{"field": field.Name}, "Skipping top-level meta field during recursion")
			continue
		}

		// Skip meta fields for child resources
		if !isTopLevel && (field.Name == "Name" || field.Name == "Namespace" || field.Name == "OwnerRef" || field.Name == "CreationTimestamp") {
			logger.Trace(logrus.Fields{"field": field.Name}, "Skipping meta field for child resource")
			continue
		}

		logger.Trace(logrus.Fields{"field": field.Name, "kind": fieldValue.Kind()}, "Processing field")

		if fieldValue.Kind() == reflect.Struct && !fieldValue.IsZero() {
			buffer.WriteString(fmt.Sprintf("%s%s:\n", indentation, field.Name))
			formatResource(buffer, fieldValue.Addr().Interface(), indent+1, false)
			logger.Trace(logrus.Fields{"field": field.Name}, "Processed struct field")
		} else if fieldValue.Kind() == reflect.Slice && fieldValue.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("%s%s:\n", indentation, field.Name))
			for j := 0; j < fieldValue.Len(); j++ {
				buffer.WriteString(fmt.Sprintf("%s- ", strings.Repeat("  ", indent+1)))
				formatResource(buffer, fieldValue.Index(j).Addr().Interface(), indent+2, false)
			}
			logger.Trace(logrus.Fields{"field": field.Name, "length": fieldValue.Len()}, "Processed slice field")
		} else if !fieldValue.IsZero() {
			buffer.WriteString(fmt.Sprintf("%s%s: %v\n", indentation, field.Name, fieldValue.Interface()))
			logger.Trace(logrus.Fields{"field": field.Name, "value": fieldValue.Interface()}, "Processed simple field")
		} else {
			logger.Trace(logrus.Fields{"field": field.Name}, "Field is zero or empty, skipping")
		}
	}

	logger.Trace(logrus.Fields{
		"resource": typ.Name(),
	}, "Completed formatResource")
}
