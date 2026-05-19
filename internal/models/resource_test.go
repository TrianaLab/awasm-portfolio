package models_test

import (
	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"testing"
	"time"
)

type mockResource struct {
	kind              string
	name              string
	namespace         string
	ownerRef          models.OwnerReference
	creationTimestamp time.Time
}

func (m *mockResource) GetKind() string {
	return m.kind
}

func (m *mockResource) GetName() string {
	return m.name
}

func (m *mockResource) SetName(name string) {
	m.name = name
}

func (m *mockResource) GetNamespace() string {
	return m.namespace
}

func (m *mockResource) SetNamespace(namespace string) {
	m.namespace = namespace
}

func (m *mockResource) GetOwnerReference() models.OwnerReference {
	return m.ownerRef
}

func (m *mockResource) SetOwnerReference(owner models.OwnerReference) {
	m.ownerRef = owner
}

func (m *mockResource) GetID() string {
	return m.kind + ":" + m.name + ":" + m.namespace
}

func (m *mockResource) GetCreationTimestamp() time.Time {
	return m.creationTimestamp
}

func (m *mockResource) SetCreationTimestamp(timestamp time.Time) {
	m.creationTimestamp = timestamp
}

func TestOwnerReferenceGetID(t *testing.T) {
	tests := []struct {
		name     string
		owner    models.OwnerReference
		expected string
	}{
		{
			name: "basic owner reference",
			owner: models.OwnerReference{
				Kind:      "Test",
				Name:      "test-name",
				Namespace: "test-ns",
			},
			expected: "test:test-name:test-ns",
		},
		{
			name: "empty namespace",
			owner: models.OwnerReference{
				Kind: "Test",
				Name: "test-name",
			},
			expected: "test:test-name:",
		},
		{
			name: "with uppercase",
			owner: models.OwnerReference{
				Kind:      "TEST",
				Name:      "TEST-name",
				Namespace: "TEST-ns",
			},
			expected: "test:test-name:test-ns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.owner.GetID()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestOwnerReferenceGetName(t *testing.T) {
	mockOwner := &mockResource{
		name: "mock-resource",
	}

	tests := []struct {
		name     string
		owner    models.OwnerReference
		expected string
	}{
		{
			name: "with owner resource",
			owner: models.OwnerReference{
				Name:  "reference-name",
				Owner: mockOwner,
			},
			expected: "mock-resource",
		},
		{
			name: "without owner resource",
			owner: models.OwnerReference{
				Name: "reference-name",
			},
			expected: "reference-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.owner.GetName()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestMeta exercises the methods promoted to concrete types via embedding.
func TestMeta(t *testing.T) {
	m := &models.Meta{}
	m.SetName("widget")
	m.SetNamespace("ns-a")
	m.SetOwnerReference(models.OwnerReference{Kind: "resume", Name: "owner"})
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetCreationTimestamp(ts)
	m.Kind = "award"

	if m.GetKind() != "award" {
		t.Errorf("GetKind() = %q, want %q", m.GetKind(), "award")
	}
	if m.GetName() != "widget" {
		t.Errorf("GetName() = %q, want %q", m.GetName(), "widget")
	}
	if m.GetNamespace() != "ns-a" {
		t.Errorf("GetNamespace() = %q, want %q", m.GetNamespace(), "ns-a")
	}
	if m.GetOwnerReference().Kind != "resume" {
		t.Errorf("GetOwnerReference().Kind = %q, want %q", m.GetOwnerReference().Kind, "resume")
	}
	if !m.GetCreationTimestamp().Equal(ts) {
		t.Errorf("GetCreationTimestamp() = %v, want %v", m.GetCreationTimestamp(), ts)
	}
	if got := m.GetID(); got != "award:widget:ns-a" {
		t.Errorf("GetID() = %q, want %q", got, "award:widget:ns-a")
	}
}

func TestMockResource(t *testing.T) {
	resource := &mockResource{
		kind:      "test",
		name:      "test-name",
		namespace: "test-ns",
	}

	t.Run("getters", func(t *testing.T) {
		if resource.GetKind() != "test" {
			t.Errorf("GetKind() = %v, want %v", resource.GetKind(), "test")
		}
		if resource.GetName() != "test-name" {
			t.Errorf("GetName() = %v, want %v", resource.GetName(), "test-name")
		}
		if resource.GetNamespace() != "test-ns" {
			t.Errorf("GetNamespace() = %v, want %v", resource.GetNamespace(), "test-ns")
		}
		if resource.GetID() != "test:test-name:test-ns" {
			t.Errorf("GetID() = %v, want %v", resource.GetID(), "test:test-name:test-ns")
		}
	})

	t.Run("setters", func(t *testing.T) {
		now := time.Now()
		resource.SetName("new-name")
		resource.SetNamespace("new-ns")
		resource.SetCreationTimestamp(now)
		resource.SetOwnerReference(models.OwnerReference{Kind: "owner", Name: "owner-name"})

		if resource.GetName() != "new-name" {
			t.Errorf("after SetName() = %v, want %v", resource.GetName(), "new-name")
		}
		if resource.GetNamespace() != "new-ns" {
			t.Errorf("after SetNamespace() = %v, want %v", resource.GetNamespace(), "new-ns")
		}
		if !resource.GetCreationTimestamp().Equal(now) {
			t.Errorf("after SetCreationTimestamp() = %v, want %v", resource.GetCreationTimestamp(), now)
		}
		if resource.GetOwnerReference().Kind != "owner" {
			t.Errorf("after SetOwnerReference() = %v, want %v", resource.GetOwnerReference().Kind, "owner")
		}
	})
}
