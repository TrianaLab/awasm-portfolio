package ui

import "awasm-portfolio/internal/models"

// Formatter defines the unified interface for formatting resources
type Formatter interface {
	FormatTable(resources []models.Resource) string
	FormatDetails(resource models.Resource) string
}

// UnifiedFormatter delegates to specific formatters for table and details
type UnifiedFormatter struct {
	tableFormatter   TableFormatter
	detailsFormatter DetailsFormatter
}

// NewUnifiedFormatter creates a new UnifiedFormatter instance
func NewUnifiedFormatter() *UnifiedFormatter {
	return &UnifiedFormatter{
		tableFormatter:   TableFormatter{},
		detailsFormatter: DetailsFormatter{},
	}
}

// FormatTable delegates to the TableFormatter
func (f *UnifiedFormatter) FormatTable(resources []models.Resource) string {
	return f.tableFormatter.FormatTable(resources)
}

// FormatDetails delegates to the DetailsFormatter
func (f *UnifiedFormatter) FormatDetails(resource models.Resource) string {
	return f.detailsFormatter.FormatDetails(resource)
}
