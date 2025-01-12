package ui

import "awasm-portfolio/internal/models"

type Formatter interface {
	FormatTable(resources []models.Resource) string
	FormatDetails(resource models.Resource) string
}
