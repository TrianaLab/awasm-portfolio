package ui

import (
	"awasm-portfolio/internal/models"
	"fmt"
	"strings"
)

type TextFormatter struct{}

func (f TextFormatter) FormatTable(data map[string]models.ResourceBase) string {
	var headers []string
	rows := [][]string{}

	// Determine if the data represents namespaces
	isNamespace := false
	for _, resource := range data {
		if resource.Namespace == "" {
			isNamespace = true
			break
		}
	}

	// Set headers based on the resource type
	if isNamespace {
		headers = []string{"Name"} // Only "Name" for namespaces
	} else {
		headers = []string{"Name", "Namespace"} // Default headers for non-namespace resources
		headerSet := map[string]bool{"Name": true, "Namespace": true}

		// Dynamically collect additional headers
		for _, resource := range data {
			fields := resource.GetFields()
			for key, val := range fields {
				if key != "Name" && key != "Namespace" && val != "" && !headerSet[key] {
					headerSet[key] = true
					headers = append(headers, strings.Title(key)) // Capitalize header
				}
			}
		}
	}

	// Collect rows
	for _, resource := range data {
		fields := resource.GetFields()
		var row []string

		if isNamespace {
			// For namespaces, only include the "Name" field
			row = []string{fields["Name"]}
		} else {
			// For non-namespace resources, ensure "Name" and "Namespace" are always included
			row = []string{
				fields["Name"],      // Include Name
				fields["Namespace"], // Include Namespace
			}

			// Include additional headers dynamically
			for _, header := range headers[2:] { // Skip Name and Namespace headers
				fieldKey := strings.ToLower(header) // Match header to field key
				if val, exists := fields[fieldKey]; exists {
					row = append(row, val)
				} else {
					row = append(row, "") // Ensure alignment for missing fields
				}
			}
		}

		rows = append(rows, row)
	}

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
	output.WriteString("\r\n") // Add newline after headers

	for _, row := range rows {
		for i, field := range row {
			output.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], field))
		}
		output.WriteString("\r\n")
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
