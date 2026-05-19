package repository_test

import (
	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"strings"
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

func TestInMemoryRepository(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	now := time.Now()

	resource := &mockResource{
		kind:              "resume",
		name:              "test-resume",
		namespace:         "default",
		creationTimestamp: now,
	}

	t.Run("Create", func(t *testing.T) {
		namespace := &mockResource{
			kind:              "namespace",
			name:              "default",
			creationTimestamp: now,
		}
		_, err := repo.Create(namespace)
		if err != nil {
			t.Fatalf("unexpected error creating namespace: %v", err)
		}

		msg, err := repo.Create(resource)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(msg, "created") {
			t.Errorf("expected created message, got: %s", msg)
		}

		_, err = repo.Create(resource)
		if err == nil {
			t.Error("expected error creating duplicate resource")
		}
	})

	t.Run("List", func(t *testing.T) {
		resources, err := repo.List(resource.GetKind(), resource.GetName(), resource.GetNamespace())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resources) != 1 {
			t.Errorf("expected 1 resource, got %d", len(resources))
		}
		if resources[0].GetID() != resource.GetID() {
			t.Errorf("expected ID %s, got %s", resource.GetID(), resources[0].GetID())
		}

		_, err = repo.List("nonexistent", "", "")
		if err == nil {
			t.Error("expected error for invalid kind")
		}

		resources, err = repo.List("resume", "", "invalid-ns")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resources) != 0 {
			t.Errorf("expected 0 resources, got %d", len(resources))
		}
	})

	t.Run("Delete", func(t *testing.T) {
		msg, err := repo.Delete(resource.GetKind(), resource.GetName(), resource.GetNamespace())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(msg, "deleted") {
			t.Errorf("expected deleted message, got: %s", msg)
		}

		_, err = repo.Delete("resume", "nonexistent", "default")
		if err == nil {
			t.Error("expected error deleting nonexistent resource")
		}

		_, err = repo.Delete("invalid", "name", "default")
		if err == nil {
			t.Error("expected error for invalid kind")
		}
	})

	t.Run("List All Namespaces", func(t *testing.T) {
		resource2 := &mockResource{
			kind:              "resume",
			name:              "test-resume-2",
			namespace:         "other",
			creationTimestamp: now,
		}

		namespace2 := &mockResource{
			kind:              "namespace",
			name:              "other",
			creationTimestamp: now,
		}
		_, _ = repo.Create(namespace2)
		_, _ = repo.Create(resource2)

		resources, err := repo.List("resume", "", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resources) != 1 {
			t.Errorf("expected 1 resource, got %d", len(resources))
		}

		resources, err = repo.List("namespace", "", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resources) != 2 {
			t.Errorf("expected 2 namespaces, got %d", len(resources))
		}
	})

	t.Run("Create Stamps Zero Timestamp", func(t *testing.T) {
		// A resource with a zero CreationTimestamp triggers the auto-stamp branch.
		r := &mockResource{
			kind:      "resume",
			name:      "no-ts",
			namespace: "default",
		}
		if _, err := repo.Create(r); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r.GetCreationTimestamp().IsZero() {
			t.Error("Create did not auto-stamp a zero CreationTimestamp")
		}
	})

	t.Run("Create Preserves Non-Zero Timestamp", func(t *testing.T) {
		fixed := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
		r := &mockResource{
			kind:              "resume",
			name:              "with-ts",
			namespace:         "default",
			creationTimestamp: fixed,
		}
		if _, err := repo.Create(r); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		stored, err := repo.List("resume", "with-ts", "default")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(stored) != 1 {
			t.Fatalf("expected 1 stored resource, got %d", len(stored))
		}
		if !stored[0].GetCreationTimestamp().Equal(fixed) {
			t.Errorf("Create overwrote a non-zero CreationTimestamp: got %v, want %v",
				stored[0].GetCreationTimestamp(), fixed)
		}
	})
}
