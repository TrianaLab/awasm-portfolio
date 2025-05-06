package util_test

import (
	"awasm-portfolio/internal/util"
	"strings"
	"testing"
)

func TestSupportedResources(t *testing.T) {
	resources := util.SupportedResources()

	tests := []struct {
		input    string
		expected string
	}{
		{"resume", "resume"},
		{"resumes", "resume"},
		{"namespace", "namespace"},
		{"ns", "namespace"},
		{"work", "work"},
		{"works", "work"},
		{"volunteer", "volunteer"},
		{"volunteers", "volunteer"},
		{"education", "education"},
		{"educations", "education"},
		{"award", "award"},
		{"awards", "award"},
		{"certificate", "certificate"},
		{"certificates", "certificate"},
		{"publication", "publication"},
		{"publications", "publication"},
		{"skill", "skill"},
		{"skills", "skill"},
		{"language", "language"},
		{"languages", "language"},
		{"interest", "interest"},
		{"interests", "interest"},
		{"reference", "reference"},
		{"references", "reference"},
		{"project", "project"},
		{"projects", "project"},
	}

	for _, tc := range tests {
		if val, exists := resources[tc.input]; !exists || val != tc.expected {
			t.Errorf("SupportedResources: expected key %q to map to %q, got %q", tc.input, tc.expected, val)
		}
	}

	if _, exists := resources["unsupportedKind"]; exists {
		t.Error("SupportedResources: unexpected key 'unsupportedKind' found")
	}
}

func TestNormalizeKind(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
		expectError    bool
		expectedErrMsg string
	}{
		// Test valid kinds
		{"resume", "resume", false, ""},
		{"Resumes", "resume", false, ""},
		{"  ns  ", "namespace", false, ""},
		{"work", "work", false, ""},
		{"Works", "work", false, ""},
		{"volunteer", "volunteer", false, ""},
		{"education", "education", false, ""},
		{"award", "award", false, ""},
		{"certificate", "certificate", false, ""},
		{"publication", "publication", false, ""},
		{"skill", "skill", false, ""},
		{"language", "language", false, ""},
		{"interest", "interest", false, ""},
		{"reference", "reference", false, ""},
		{"project", "project", false, ""},
		{"All", "", false, ""}, // caso especial para "all"
		// Test invalid kinds
		{"invalidKind", "", true, "unsupported resource kind: invalidkind"},
		{"", "", true, "unsupported resource kind: "}, // cadena vacía debe ser no soportada
	}

	for _, tc := range tests {
		output, err := util.NormalizeKind(tc.input)
		if tc.expectError {
			if err == nil {
				t.Errorf("NormalizeKind(%q): expected error but got none", tc.input)
			} else if !strings.Contains(err.Error(), tc.expectedErrMsg) {
				t.Errorf("NormalizeKind(%q): expected error containing %q, got %q", tc.input, tc.expectedErrMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("NormalizeKind(%q): unexpected error: %v", tc.input, err)
			}
			if output != tc.expectedOutput {
				t.Errorf("NormalizeKind(%q): expected %q, got %q", tc.input, tc.expectedOutput, output)
			}
		}
	}
}
