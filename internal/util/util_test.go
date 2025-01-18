package util_test

import (
	"awasm-portfolio/internal/util"
	"strings"
	"testing"
)

func TestSupportedResources(t *testing.T) {
	resources := util.SupportedResources()

	// Check that certain keys map to expected canonical values.
	tests := []struct {
		input    string
		expected string
	}{
		{"profile", "profile"},
		{"profiles", "profile"},
		{"namespace", "namespace"},
		{"ns", "namespace"},
		{"education", "education"},
		{"experience", "experience"},
		{"contacts", "contact"},
		{"certification", "certifications"},
		{"contributions", "contributions"},
		{"skill", "skills"},
	}

	for _, tc := range tests {
		if val, exists := resources[tc.input]; !exists || val != tc.expected {
			t.Errorf("SupportedResources: expected key %q to map to %q, got %q", tc.input, tc.expected, val)
		}
	}

	// Check that unsupported keys are not present.
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
		{"profile", "profile", false, ""},
		{"Profiles", "profile", false, ""},
		{"  ns  ", "namespace", false, ""},
		{"All", "", false, ""}, // special case for "all"
		// Test invalid kind
		{"invalidKind", "", true, "unsupported resource kind: invalidkind"},
		{"", "", true, "unsupported resource kind: "}, // empty string should be unsupported
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
