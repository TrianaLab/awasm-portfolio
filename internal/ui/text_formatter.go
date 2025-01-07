package ui

import (
	"fmt"
	"sort"
	"strings"
)

type TextFormatter struct{}

func (f TextFormatter) FormatTable(data map[string]interface{}) string {
	headers := []string{}
	rows := [][]string{}

	// Collect headers and rows dynamically
	for _, value := range data {
		resource, ok := value.(interface{ GetFields() map[string]string })
		if !ok {
			continue
		}

		fields := resource.GetFields()
		if len(headers) == 0 {
			for field := range fields {
				headers = append(headers, field)
			}
			sort.Strings(headers) // Ensure consistent order
		}

		row := []string{}
		for _, header := range headers {
			row = append(row, fields[header])
		}
		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	// Determine column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, field := range row {
			if len(field) > colWidths[i] {
				colWidths[i] = len(field)
			}
		}
	}

	// Build the table output
	var output strings.Builder
	for i, header := range headers {
		output.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], header))
	}
	output.WriteString("\n")

	for _, row := range rows {
		for i, field := range row {
			output.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], field))
		}
		output.WriteString("\n")
	}

	return strings.TrimSuffix(output.String(), "\n")
}

func (f TextFormatter) FormatDetails(data map[string]string) string {
	var output strings.Builder
	for key, value := range data {
		output.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	return strings.TrimSuffix(output.String(), "\n")
}

func (f TextFormatter) Error(message string) string {
	return fmt.Sprintf("Error: %s", message)
}

func (f TextFormatter) Success(message string) string {
	return fmt.Sprintf("Success: %s", message)
}
