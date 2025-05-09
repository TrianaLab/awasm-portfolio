package models_test

import (
	"awasm-portfolio/internal/models"
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
