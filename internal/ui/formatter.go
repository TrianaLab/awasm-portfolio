package ui

import "awasm-portfolio/internal/models"

type Formatter interface {
	FormatTable(data map[string]models.ResourceBase) string
	FormatDetails(data map[string]string) string
	Error(message string) string
	Success(message string) string
}
