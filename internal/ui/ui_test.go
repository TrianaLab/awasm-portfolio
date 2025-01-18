package ui

import (
	"awasm-portfolio/internal/models"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// testResource is an implementation of models.Resource with exported fields.
type testResource struct {
	Kind              string                `yaml:"kind"`
	Name              string                `yaml:"name"`
	Namespace         string                `yaml:"namespace"`
	ID                string                `yaml:"id"`
	CreationTimestamp time.Time             `yaml:"creationTimestamp"`
	OwnerReference    models.OwnerReference `yaml:"ownerReference"`
}

func (r *testResource) GetKind() string                          { return r.Kind }
func (r *testResource) GetName() string                          { return r.Name }
func (r *testResource) SetName(name string)                      { r.Name = name }
func (r *testResource) GetNamespace() string                     { return r.Namespace }
func (r *testResource) SetNamespace(namespace string)            { r.Namespace = namespace }
func (r *testResource) GetOwnerReference() models.OwnerReference { return r.OwnerReference }
func (r *testResource) SetOwnerReference(owner models.OwnerReference) {
	r.OwnerReference = owner
}
func (r *testResource) GetID() string                    { return r.ID }
func (r *testResource) GetCreationTimestamp() time.Time  { return r.CreationTimestamp }
func (r *testResource) SetCreationTimestamp(t time.Time) { r.CreationTimestamp = t }

// newTestResource is a helper to create a new testResource.
func newTestResource(kind, name, namespace string) *testResource {
	return &testResource{
		Kind:              kind,
		Name:              name,
		Namespace:         namespace,
		ID:                kind + ":" + namespace + ":" + name,
		CreationTimestamp: time.Now().Add(-2 * time.Hour), // fixed relative timestamp
	}
}

func TestFormatDetails(t *testing.T) {
	res := newTestResource("testKind", "testName", "testNS")
	resources := []models.Resource{res}

	output := FormatDetails(resources)

	var parsed []testResource
	err := yaml.Unmarshal([]byte(output), &parsed)
	if err != nil {
		t.Fatalf("failed to unmarshal YAML: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected one resource in YAML output, got %d", len(parsed))
	}

	unmarshaled := parsed[0]

	// Compare fields of the original resource and the unmarshaled resource.
	if unmarshaled.Kind != res.Kind {
		t.Errorf("expected kind %q, got %q", res.Kind, unmarshaled.Kind)
	}
	if unmarshaled.Name != res.Name {
		t.Errorf("expected name %q, got %q", res.Name, unmarshaled.Name)
	}
	if unmarshaled.Namespace != res.Namespace {
		t.Errorf("expected namespace %q, got %q", res.Namespace, unmarshaled.Namespace)
	}
	if !unmarshaled.CreationTimestamp.Equal(res.CreationTimestamp) {
		t.Errorf("expected creationTimestamp %v, got %v", res.CreationTimestamp, unmarshaled.CreationTimestamp)
	}
}

func TestFormatTableJSON(t *testing.T) {
	res := newTestResource("testKind", "testName", "testNS")
	resources := []models.Resource{res}

	output := FormatTable(resources, "json")

	// Try to unmarshal JSON to ensure it's valid.
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &data); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(data) != 1 {
		t.Errorf("expected 1 JSON object, got %d", len(data))
	}
}

func TestFormatTableYAML(t *testing.T) {
	res := newTestResource("testKind", "testName", "testNS")
	resources := []models.Resource{res}

	output := FormatTable(resources, "yaml")

	// Try to unmarshal YAML to ensure it's valid.
	var parsed []map[string]interface{}
	if err := yaml.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid YAML: %v", err)
	}
	if len(parsed) != 1 {
		t.Errorf("expected 1 YAML document, got %d", len(parsed))
	}
}

func TestFormatTableDefault(t *testing.T) {
	// Create a namespace resource to trigger the "namespace" schema.
	res := newTestResource("namespace", "ns1", "")
	resources := []models.Resource{res}

	// Use default table format
	output := FormatTable(resources, "table")

	// Check that output contains header keywords for namespace.
	if !strings.Contains(output, "NAME") || !strings.Contains(output, "AGE") {
		t.Errorf("table output missing expected headers, got: %s", output)
	}

	// Also check that resource name appears in output.
	if !strings.Contains(output, "ns1") {
		t.Errorf("table output missing resource name 'ns1', got: %s", output)
	}
}
