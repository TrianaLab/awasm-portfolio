package ui

type Formatter interface {
	FormatTable(data map[string]interface{}) string
	FormatDetails(data map[string]string) string
	Error(message string) string
	Success(message string) string
}
