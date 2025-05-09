package ui_test

import (
	"awasm-portfolio/internal/models"
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
