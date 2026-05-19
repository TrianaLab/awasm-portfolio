package ui_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
	"github.com/TrianaLab/awasm-portfolio/internal/ui"

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

func meta(kind, name, namespace string, ts time.Time) models.Meta {
	return models.Meta{Kind: kind, Name: name, Namespace: namespace, CreationTimestamp: ts}
}

func TestFormatTable_Table(t *testing.T) {
	now := time.Now()
	res := &testResource{Kind: "resume", Name: "test", Namespace: "default", CreationTimestamp: now}
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
	res := &testResource{Kind: "resume", Name: "test", Namespace: "default", CreationTimestamp: now}
	out := ui.FormatTable([]models.Resource{res}, "json")
	var v []map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Errorf("FormatTable json output is not valid JSON: %v", err)
	}
	if empty := ui.FormatTable([]models.Resource{}, "json"); empty != "[]" {
		t.Errorf("FormatTable json output for empty slice should be []: %s", empty)
	}
	if nilOut := ui.FormatTable(nil, "json"); nilOut != "null" {
		t.Errorf("FormatTable json output for nil should be null: %s", nilOut)
	}
}

func TestFormatTable_YAML(t *testing.T) {
	now := time.Now()
	res := &testResource{Kind: "resume", Name: "test", Namespace: "default", CreationTimestamp: now}
	out := ui.FormatTable([]models.Resource{res}, "yaml")
	var v []map[string]any
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
	res := &testResource{Kind: "resume", Name: "test", Namespace: "default", CreationTimestamp: time.Now()}
	out := ui.FormatTable([]models.Resource{res}, "unknown")
	if !strings.Contains(out, "NAME") {
		t.Error("FormatTable with unknown format should fallback to table")
	}
}

func TestFormatTable_CaseInsensitiveFormat(t *testing.T) {
	res := &testResource{Kind: "resume", Name: "test", Namespace: "default", CreationTimestamp: time.Now()}

	out := ui.FormatTable([]models.Resource{res}, "JSON")
	var v []map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Errorf("FormatTable json (case-insensitive) output is not valid JSON: %v", err)
	}

	out = ui.FormatTable([]models.Resource{res}, "YAML")
	var v2 []map[string]any
	if err := yaml.Unmarshal([]byte(out), &v2); err != nil {
		t.Errorf("FormatTable yaml (case-insensitive) output is not valid YAML: %v", err)
	}
}

func TestFormatDetails(t *testing.T) {
	res := &testResource{
		Kind: "resume", Name: "test", Namespace: "default", CreationTimestamp: time.Now(),
		ownerRef: models.OwnerReference{Kind: "parent", Name: "parentName", Namespace: "parentNS"},
	}
	out := ui.FormatDetails([]models.Resource{res})
	if !strings.Contains(out, "kind: resume") || !strings.Contains(out, "name: test") {
		t.Errorf("FormatDetails output missing expected content: %s", out)
	}
	if empty := ui.FormatDetails([]models.Resource{}); !strings.Contains(empty, "[]") {
		t.Error("FormatDetails with empty slice should return empty array")
	}
	if nilOut := ui.FormatDetails(nil); !strings.Contains(nilOut, "[]") {
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
	schema := schemas["namespace"]
	resource := &types.Namespace{Meta: meta("namespace", "test-namespace", "", time.Now())}
	for i, extractor := range schema.Extractors {
		if extractor(resource) == "" {
			t.Errorf("Extractor %d for namespace returned empty value", i)
		}
	}
}

func TestDefaultSchema(t *testing.T) {
	schemas := ui.GenerateSchemas()
	schema := schemas["default"]
	resource := &testResource{Kind: "default", Name: "test-default", Namespace: "default-ns", CreationTimestamp: time.Now()}
	for i, extractor := range schema.Extractors {
		if extractor(resource) == "" {
			t.Errorf("Extractor %d for default schema returned empty value", i)
		}
	}
}

func TestResumeSchema(t *testing.T) {
	schemas := ui.GenerateSchemas()
	schema := schemas["resume"]
	resource := &types.Resume{Meta: meta("resume", "test-resume", "default-ns", time.Now())}
	for _, extractor := range schema.Extractors {
		_ = extractor(resource)
	}
}

func TestAllSchemas(t *testing.T) {
	schemas := ui.GenerateSchemas()
	now := time.Now()

	resumeRes := &types.Resume{
		Meta:         meta("resume", "test-resume", "default-ns", now),
		Basics:       types.Basics{Meta: meta("basics", "John Doe", "", time.Time{})},
		Work:         []types.Work{{Company: "Test Company"}},
		Volunteer:    []types.Volunteer{{Organization: "Test Org"}},
		Education:    []types.Education{{Institution: "Test University"}},
		Awards:       []types.Award{{Title: "Best Developer"}},
		Certificates: []types.Certificate{{Certificate: "Test Certificate"}},
		Publications: []types.Publication{{Publication: "Test Publication"}},
		Skills:       []types.Skill{{Skill: "Go Programming"}},
		Languages:    []types.Language{{Language: "English"}},
		Interests:    []types.Interest{{Interest: "Programming"}},
		References:   []types.Reference{{Meta: meta("reference", "Test Reference", "", time.Time{})}},
		Projects:     []types.Project{{Project: "Test Project"}},
	}

	testCases := map[string]models.Resource{
		"namespace":   &types.Namespace{Meta: meta("namespace", "test-namespace", "", now)},
		"resume":      resumeRes,
		"basics":      &types.Basics{Meta: meta("basics", "test-basics", "default-ns", now), FullName: "John Doe", Label: "Engineer", Email: "john.doe@example.com", Phone: "123456789"},
		"work":        &types.Work{Meta: meta("work", "test-work", "default-ns", now), Company: "Test Company", Position: "Developer", StartDate: "2020-01-01", EndDate: "2022-01-01"},
		"volunteer":   &types.Volunteer{Meta: meta("volunteer", "test-volunteer", "default-ns", now), Organization: "Test Org", Position: "Volunteer"},
		"education":   &types.Education{Meta: meta("education", "test-education", "default-ns", now), Institution: "Test University", Area: "Computer Science", StudyType: "Bachelor"},
		"skill":       &types.Skill{Meta: meta("skill", "test-skill", "default-ns", now), Skill: "Go Programming", Level: "Expert", Keywords: []string{"Go", "Programming"}},
		"language":    &types.Language{Meta: meta("language", "test-language", "default-ns", now), Language: "English", Fluency: "Native"},
		"project":     &types.Project{Meta: meta("project", "test-project", "default-ns", now), Project: "Test Project", URL: "https://example.com"},
		"publication": &types.Publication{Meta: meta("publication", "test-publication", "default-ns", now), Publication: "Test Publication", Publisher: "Test Publisher"},
		"certificate": &types.Certificate{Meta: meta("certificate", "test-certificate", "default-ns", now), Certificate: "Test Certificate", Date: "2022-01-01", Issuer: "Test Issuer"},
		"interest":    &types.Interest{Meta: meta("interest", "test-interest", "default-ns", now), Interest: "Programming"},
		"award":       &types.Award{Meta: meta("award", "test-award", "default-ns", now), Title: "Best Developer", Awarder: "Test Org", Date: "2022-01-01"},
		"default":     &testResource{Kind: "default", Name: "test-default", Namespace: "default-ns", CreationTimestamp: now},
	}

	for kind, resource := range testCases {
		t.Run(kind, func(t *testing.T) {
			schema, ok := schemas[kind]
			if !ok {
				t.Fatalf("schema for %s not found", kind)
			}
			for i, extractor := range schema.Extractors {
				if extractor(resource) == "" {
					t.Errorf("Extractor %d for schema %s returned empty value", i, kind)
				}
			}
		})
	}
}

// TestSchemaTypeAssertionFallback feeds a wrong-type resource into every
// typed schema's extractors to exercise the "N/A" fallback paths that fire
// when the type assertion in each extractor fails.
func TestSchemaTypeAssertionFallback(t *testing.T) {
	schemas := ui.GenerateSchemas()
	other := &testResource{Kind: "other", Name: "wrong", Namespace: "default", CreationTimestamp: time.Now()}

	for kind, schema := range schemas {
		if kind == "namespace" || kind == "default" {
			continue // no type assertions in these schemas
		}
		t.Run(kind, func(t *testing.T) {
			for _, extractor := range schema.Extractors {
				_ = extractor(other)
			}
		})
	}
}

func TestWorkSchema_PresentEndDate(t *testing.T) {
	schemas := ui.GenerateSchemas()
	work := &types.Work{
		Meta:      meta("work", "current-role", "default", time.Now()),
		Company:   "Acme",
		Position:  "Engineer",
		StartDate: "2024-01-01",
	}
	found := false
	for _, extractor := range schemas["work"].Extractors {
		if extractor(work) == "Present" {
			found = true
			break
		}
	}
	if !found {
		t.Error("work schema should render \"Present\" when EndDate is empty")
	}
}

func TestSkillSchema_KeywordsTruncation(t *testing.T) {
	schemas := ui.GenerateSchemas()
	skill := &types.Skill{
		Meta:     meta("skill", "many", "default", time.Now()),
		Skill:    "Polyglot",
		Level:    "Expert",
		Keywords: []string{"Go", "Rust", "Python", "Zig", "C"},
	}
	found := false
	for _, extractor := range schemas["skill"].Extractors {
		if v := extractor(skill); strings.Contains(v, "...") {
			found = true
			break
		}
	}
	if !found {
		t.Error("skill schema should truncate keyword lists with > 3 entries")
	}
}

func TestCalculateAge(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name   string
		ts     time.Time
		expect string
	}{
		{"zero", time.Time{}, ""},
		{"seconds", now.Add(-5 * time.Second), "s"},
		{"minutes", now.Add(-5 * time.Minute), "m"},
		{"hours", now.Add(-3 * time.Hour), "h"},
		{"days", now.Add(-3 * 24 * time.Hour), "d"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := &testResource{Kind: "namespace", Name: "ns-" + tc.name, CreationTimestamp: tc.ts}
			out := ui.FormatTable([]models.Resource{res}, "")
			if tc.expect != "" && !strings.Contains(out, tc.expect) {
				t.Errorf("expected %q in age output, got: %s", tc.expect, out)
			}
		})
	}
}

func TestFormatAsTable_UnknownKindFallback(t *testing.T) {
	res := &testResource{Kind: "definitely-not-a-real-kind", Name: "weird", Namespace: "default", CreationTimestamp: time.Now()}
	out := ui.FormatTable([]models.Resource{res}, "")
	if !strings.Contains(out, "weird") {
		t.Errorf("default-schema fallback should render the resource, got: %s", out)
	}
}

func TestFormatAsTable_WideCell(t *testing.T) {
	long := strings.Repeat("x", 80)
	res := &testResource{Kind: "default", Name: long, Namespace: "default", CreationTimestamp: time.Now()}
	out := ui.FormatTable([]models.Resource{res}, "")
	if !strings.Contains(out, long) {
		t.Errorf("wide cell should render in full, got: %s", out)
	}
}

// failingResource implements models.Resource but errors on JSON/YAML marshal,
// exercising the marshal-error fallback paths.
type failingResource struct{}

func (f *failingResource) GetKind() string                          { return "failing" }
func (f *failingResource) GetName() string                          { return "fail" }
func (f *failingResource) SetName(string)                           {}
func (f *failingResource) GetNamespace() string                     { return "default" }
func (f *failingResource) SetNamespace(string)                      {}
func (f *failingResource) GetOwnerReference() models.OwnerReference { return models.OwnerReference{} }
func (f *failingResource) SetOwnerReference(models.OwnerReference)  {}
func (f *failingResource) GetID() string                            { return "failing:fail:default" }
func (f *failingResource) GetCreationTimestamp() time.Time          { return time.Time{} }
func (f *failingResource) SetCreationTimestamp(time.Time)           {}

func (f *failingResource) MarshalJSON() ([]byte, error) {
	return nil, errors.New("intentional json marshal failure")
}

func (f *failingResource) MarshalYAML() (any, error) {
	return nil, errors.New("intentional yaml marshal failure")
}

func TestFormatAsJSON_Error(t *testing.T) {
	out := ui.FormatTable([]models.Resource{&failingResource{}}, "json")
	if !strings.Contains(out, "Error formatting resources as JSON") {
		t.Errorf("expected JSON error message, got: %s", out)
	}
}

func TestFormatAsYAML_Error(t *testing.T) {
	out := ui.FormatTable([]models.Resource{&failingResource{}}, "yaml")
	if !strings.Contains(out, "Error formatting resources as YAML") {
		t.Errorf("expected YAML error message, got: %s", out)
	}
}
