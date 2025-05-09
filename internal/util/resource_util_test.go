package util_test

import (
	"awasm-portfolio/internal/util"
	"testing"
)

func TestSupportedResources(t *testing.T) {
	resources := util.SupportedResources()

	// Verificar que algunos recursos clave están presentes
	expectedResources := map[string]string{
		"resume":      "resume",
		"resumes":     "resume",
		"basic":       "basics",
		"basics":      "basics",
		"namespace":   "namespace",
		"namespaces":  "namespace",
		"work":        "work",
		"volunteer":   "volunteer",
		"education":   "education",
		"award":       "award",
		"certificate": "certificate",
		"publication": "publication",
		"skill":       "skill",
		"language":    "language",
		"interest":    "interest",
		"reference":   "reference",
		"project":     "project",
	}

	for key, expected := range expectedResources {
		if resources[key] != expected {
			t.Errorf("Expected %s to map to %s, but got %s", key, expected, resources[key])
		}
	}

	// Verificar que el mapa no está vacío
	if len(resources) == 0 {
		t.Error("SupportedResources returned an empty map")
	}
}

func TestNormalizeKind(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"resume", "resume", false},
		{"resumes", "resume", false},
		{"basic", "basics", false},
		{"basics", "basics", false},
		{"namespace", "namespace", false},
		{"namespaces", "namespace", false},
		{"work", "work", false},
		{"volunteer", "volunteer", false},
		{"education", "education", false},
		{"award", "award", false},
		{"certificate", "certificate", false},
		{"publication", "publication", false},
		{"skill", "skill", false},
		{"language", "language", false},
		{"interest", "interest", false},
		{"reference", "reference", false},
		{"project", "project", false},
		{"all", "", false},                // Caso especial para "all"
		{"unknown", "", true},             // Caso de error
		{"", "", true},                    // Caso vacío
		{"   resume   ", "resume", false}, // Caso con espacios
		{"RESUME", "resume", false},       // Caso con mayúsculas
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := util.NormalizeKind(test.input)
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error for input %s, but got: %v", test.input, err)
				}
				if result != test.expected {
					t.Errorf("Expected %s for input %s, but got %s", test.expected, test.input, result)
				}
			}
		})
	}
}
