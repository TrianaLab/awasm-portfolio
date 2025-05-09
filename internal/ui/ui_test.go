package ui_test

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/ui"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

type testResource struct {
	Kind              string
	Name              string
	Namespace         string
	CreationTimestamp time.Time
	ownerRef          models.OwnerReference
}

func (t *testResource) GetKind() string                             { return t.Kind }
func (t *testResource) GetName() string                             { return t.Name }
func (t *testResource) SetName(name string)                         { t.Name = name }
func (t *testResource) GetNamespace() string                        { return t.Namespace }
func (t *testResource) SetNamespace(namespace string)               { t.Namespace = namespace }
func (t *testResource) GetOwnerReference() models.OwnerReference    { return t.ownerRef }
func (t *testResource) SetOwnerReference(ref models.OwnerReference) { t.ownerRef = ref }
func (t *testResource) GetID() string                               { return t.Kind + ":" + t.Name + ":" + t.Namespace }
func (t *testResource) GetCreationTimestamp() time.Time             { return t.CreationTimestamp }
func (t *testResource) SetCreationTimestamp(timestamp time.Time)    { t.CreationTimestamp = timestamp }

func TestFormatTable_Table(t *testing.T) {
	now := time.Now()
	res := &testResource{
		Kind:              "resume",
		Name:              "test",
		Namespace:         "default",
		CreationTimestamp: now,
	}
	out := ui.FormatTable([]models.Resource{res}, "")
	if !strings.Contains(out, "NAME") || !strings.Contains(out, "test") {
		t.Errorf("FormatTable table output missing expected content: %s", out)
	}
	empty := ui.FormatTable([]models.Resource{}, "")
	if !strings.Contains(empty, "No resources found") {
		t.Error("FormatTable with empty slice should show 'No resources found'")
	}
	nilOut := ui.FormatTable(nil, "")
	if !strings.Contains(nilOut, "No resources found") {
		t.Error("FormatTable with empty slice should show 'No resources found'")
	}
}

func TestFormatTable_JSON(t *testing.T) {
	now := time.Now()
	res := &testResource{
		Kind:              "resume",
		Name:              "test",
		Namespace:         "default",
		CreationTimestamp: now,
	}
	out := ui.FormatTable([]models.Resource{res}, "json")
	var v []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Errorf("FormatTable json output is not valid JSON: %v", err)
	}
	empty := ui.FormatTable([]models.Resource{}, "json")
	if empty != "[]" {
		t.Errorf("FormatTable json output for empty slice should be []: %s", empty)
	}
	nilOut := ui.FormatTable(nil, "json")
	if nilOut != "null" {
		t.Errorf("FormatTable json output for nil should be null: %s", nilOut)
	}
}

func TestFormatTable_YAML(t *testing.T) {
	now := time.Now()
	res := &testResource{
		Kind:              "resume",
		Name:              "test",
		Namespace:         "default",
		CreationTimestamp: now,
	}
	out := ui.FormatTable([]models.Resource{res}, "yaml")
	var v []map[string]interface{}
	if err := yaml.Unmarshal([]byte(out), &v); err != nil {
		t.Errorf("FormatTable yaml output is not valid YAML: %v", err)
	}
	empty := ui.FormatTable([]models.Resource{}, "yaml")
	if !strings.Contains(empty, "[]") && strings.TrimSpace(empty) != "" {
		t.Errorf("FormatTable yaml output for empty slice should be [] or empty: %s", empty)
	}
	nilOut := ui.FormatTable(nil, "yaml")
	if !strings.Contains(nilOut, "[]") && strings.TrimSpace(nilOut) != "" {
		t.Errorf("FormatTable yaml output for nil should be [] or empty: %s", nilOut)
	}
}

func TestFormatTable_UnknownFormat(t *testing.T) {
	now := time.Now()
	res := &testResource{
		Kind:              "resume",
		Name:              "test",
		Namespace:         "default",
		CreationTimestamp: now,
	}
	out := ui.FormatTable([]models.Resource{res}, "unknown")
	if !strings.Contains(out, "NAME") {
		t.Error("FormatTable with unknown format should fallback to table")
	}
}

func TestFormatTable_CaseInsensitiveFormat(t *testing.T) {
	now := time.Now()
	res := &testResource{
		Kind:              "resume",
		Name:              "test",
		Namespace:         "default",
		CreationTimestamp: now,
	}
	out := ui.FormatTable([]models.Resource{res}, "JSON")
	var v []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Errorf("FormatTable json (case-insensitive) output is not valid JSON: %v", err)
	}
	out = ui.FormatTable([]models.Resource{res}, "YAML")
	var v2 []map[string]interface{}
	if err := yaml.Unmarshal([]byte(out), &v2); err != nil {
		t.Errorf("FormatTable yaml (case-insensitive) output is not valid YAML: %v", err)
	}
}

func TestFormatDetails(t *testing.T) {
	now := time.Now()
	res := &testResource{
		Kind:              "resume",
		Name:              "test",
		Namespace:         "default",
		CreationTimestamp: now,
		ownerRef: models.OwnerReference{
			Kind:      "parent",
			Name:      "parentName",
			Namespace: "parentNS",
		},
	}
	out := ui.FormatDetails([]models.Resource{res})
	if !strings.Contains(out, "kind: resume") || !strings.Contains(out, "name: test") {
		t.Errorf("FormatDetails output missing expected content: %s", out)
	}
	empty := ui.FormatDetails([]models.Resource{})
	if !strings.Contains(empty, "[]") {
		t.Error("FormatDetails with empty slice should return empty array")
	}
	nilOut := ui.FormatDetails(nil)
	if !strings.Contains(nilOut, "[]") {
		t.Error("FormatDetails with empty slice should return empty array")
	}
}

func TestGenerateSchemas(t *testing.T) {
	schemas := ui.GenerateSchemas()

	for kind, schema := range schemas {
		if len(schema.Headers) == 0 {
			t.Errorf("Schema %s has no headers", kind)
		}
		if len(schema.Extractors) == 0 {
			t.Errorf("Schema %s has no extractors", kind)
		}
		if len(schema.Headers) != len(schema.Extractors) {
			t.Errorf("Schema %s has mismatched headers and extractors", kind)
		}
	}
}

func TestNamespaceSchema(t *testing.T) {
	schemas := ui.GenerateSchemas()
	schema, ok := schemas["namespace"]
	if !ok {
		t.Fatal("namespace schema not found")
	}

	now := time.Now()
	resource := &types.Namespace{
		Name:              "test-namespace",
		CreationTimestamp: now,
	}

	for i, extractor := range schema.Extractors {
		value := extractor(resource)
		if value == "" {
			t.Errorf("Extractor %d for namespace returned empty value", i)
		}
	}
}

func TestDefaultSchema(t *testing.T) {
	schemas := ui.GenerateSchemas()
	schema, ok := schemas["default"]
	if !ok {
		t.Fatal("default schema not found")
	}

	now := time.Now()
	resource := &testResource{
		Kind:              "default",
		Name:              "test-default",
		Namespace:         "default-ns",
		CreationTimestamp: now,
	}

	for i, extractor := range schema.Extractors {
		value := extractor(resource)
		if value == "" {
			t.Errorf("Extractor %d for default schema returned empty value", i)
		}
	}
}

func TestResumeSchema(t *testing.T) {
	schemas := ui.GenerateSchemas()
	schema, ok := schemas["resume"]
	if !ok {
		t.Fatal("resume schema not found")
	}

	now := time.Now()
	resource := &types.Resume{
		Name:              "test-resume",
		Namespace:         "default-ns",
		CreationTimestamp: now,
	}

	for _, extractor := range schema.Extractors {
		_ = extractor(resource)
	}
}

func TestAllSchemas(t *testing.T) {
	schemas := ui.GenerateSchemas()

	now := time.Now()
	testCases := map[string]models.Resource{
		"namespace": &types.Namespace{
			Name:              "test-namespace",
			CreationTimestamp: now,
		},
		"resume": &types.Resume{
			Name:              "test-resume",
			Namespace:         "default-ns",
			Basics:            types.Basics{Name: "John Doe"},
			Work:              []types.Work{{Company: "Test Company"}},
			Volunteer:         []types.Volunteer{{Organization: "Test Org"}},
			Education:         []types.Education{{Institution: "Test University"}},
			Awards:            []types.Award{{Title: "Best Developer"}},
			Certificates:      []types.Certificate{{Certificate: "Test Certificate"}},
			Publications:      []types.Publication{{Publication: "Test Publication"}},
			Skills:            []types.Skill{{Skill: "Go Programming"}},
			Languages:         []types.Language{{Language: "English"}},
			Interests:         []types.Interest{{Interest: "Programming"}},
			References:        []types.Reference{{Name: "Test Reference"}},
			Projects:          []types.Project{{Project: "Test Project"}},
			CreationTimestamp: now,
		},
		"basics": &types.Basics{
			Name:              "test-basics",
			Namespace:         "default-ns",
			FullName:          "John Doe",
			Label:             "Engineer",
			Email:             "john.doe@example.com",
			Phone:             "123456789",
			CreationTimestamp: now,
		},
		"work": &types.Work{
			Name:              "test-work",
			Namespace:         "default-ns",
			Company:           "Test Company",
			Position:          "Developer",
			StartDate:         "2020-01-01",
			EndDate:           "2022-01-01",
			CreationTimestamp: now,
		},
		"volunteer": &types.Volunteer{
			Name:              "test-volunteer",
			Namespace:         "default-ns",
			Organization:      "Test Org",
			Position:          "Volunteer",
			CreationTimestamp: now,
		},
		"education": &types.Education{
			Name:              "test-education",
			Namespace:         "default-ns",
			Institution:       "Test University",
			Area:              "Computer Science",
			StudyType:         "Bachelor",
			CreationTimestamp: now,
		},
		"skill": &types.Skill{
			Name:              "test-skill",
			Namespace:         "default-ns",
			Skill:             "Go Programming",
			Level:             "Expert",
			Keywords:          []string{"Go", "Programming"},
			CreationTimestamp: now,
		},
		"language": &types.Language{
			Name:              "test-language",
			Namespace:         "default-ns",
			Language:          "English",
			Fluency:           "Native",
			CreationTimestamp: now,
		},
		"project": &types.Project{
			Name:              "test-project",
			Namespace:         "default-ns",
			Project:           "Test Project",
			URL:               "https://example.com",
			CreationTimestamp: now,
		},
		"publication": &types.Publication{
			Name:              "test-publication",
			Namespace:         "default-ns",
			Publication:       "Test Publication",
			Publisher:         "Test Publisher",
			CreationTimestamp: now,
		},
		"certificate": &types.Certificate{
			Name:              "test-certificate",
			Namespace:         "default-ns",
			Certificate:       "Test Certificate",
			Date:              "2022-01-01",
			Issuer:            "Test Issuer",
			CreationTimestamp: now,
		},
		"interest": &types.Interest{
			Name:              "test-interest",
			Namespace:         "default-ns",
			Interest:          "Programming",
			CreationTimestamp: now,
		},
		"award": &types.Award{
			Name:              "test-award",
			Namespace:         "default-ns",
			Title:             "Best Developer",
			Awarder:           "Test Org",
			Date:              "2022-01-01",
			CreationTimestamp: now,
		},
		"default": &testResource{
			Kind:              "default",
			Name:              "test-default",
			Namespace:         "default-ns",
			CreationTimestamp: now,
		},
	}

	for kind, resource := range testCases {
		t.Run(kind, func(t *testing.T) {
			schema, ok := schemas[kind]
			if !ok {
				t.Fatalf("schema for %s not found", kind)
			}

			for i, extractor := range schema.Extractors {
				value := extractor(resource)
				if value == "" {
					t.Errorf("Extractor %d for schema %s returned empty value", i, kind)
				}
			}
		})
	}
}
