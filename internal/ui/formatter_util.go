package ui

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

// capitalizeFieldName capitalizes the first letter of a field name
func capitalizeFieldName(fieldName string) string {
	return strings.ToUpper(fieldName[:1]) + strings.ToLower(fieldName[1:])
}

// calculateColumnWidths calculates the width for each column
func calculateColumnWidths(headers []string, rows [][]string) []int {
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
	return colWidths
}

// formatHeaders formats the table headers
func formatHeaders(sb *strings.Builder, headers []string, colWidths []int) {
	for i, header := range headers {
		sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], header))
		if i < len(headers)-1 {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("\n")
}

// formatRows formats the table rows
func formatRows(sb *strings.Builder, rows [][]string, colWidths []int) {
	for _, row := range rows {
		for i, cell := range row {
			sb.WriteString(fmt.Sprintf("%-*s", colWidths[i], cell))
			if i < len(row)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}
}

// groupResourcesByKind groups resources by their kind
func groupResourcesByKind(resources []models.Resource) map[string][]models.Resource {
	grouped := make(map[string][]models.Resource)
	for _, resource := range resources {
		kind := resource.GetKind()
		grouped[kind] = append(grouped[kind], resource)
	}
	return grouped
}

func formatOwnerRef(fieldValue reflect.Value) string {
	if fieldValue.Kind() != reflect.Struct {
		logger.Warn(logrus.Fields{
			"fieldKind": fieldValue.Kind(),
		}, "OwnerRef is not a struct")
		return ""
	}

	ownerRef, ok := fieldValue.Interface().(models.OwnerReference)
	if !ok {
		logger.Warn(logrus.Fields{
			"fieldType": fieldValue.Type().String(),
		}, "Failed to cast OwnerRef to models.OwnerReference")
		return ""
	}

	// Format as "kind/name"
	return fmt.Sprintf("%s/%s", ownerRef.Kind, ownerRef.Name)
}

func summarizeArray(fieldValue reflect.Value, header string) string {
	count := fieldValue.Len()
	if count == 0 {
		return ""
	}

	// Use the header name to create a meaningful summary
	label := strings.ToLower(strings.TrimSuffix(header, "S")) // Singular form of header
	if count > 1 {
		label += "s" // Pluralize for counts > 1
	}

	return fmt.Sprintf("%d %s", count, label)
}
