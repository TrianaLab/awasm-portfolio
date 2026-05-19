package repository_test

import (
	"strings"
	"testing"
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
)

type mockResource struct {
	kind              string
	name              string
	namespace         string
	ownerRef          models.OwnerReference
	creationTimestamp time.Time
}

func (m *mockResource) GetKind() string                          { return m.kind }
func (m *mockResource) GetName() string                          { return m.name }
func (m *mockResource) SetName(name string)                      { m.name = name }
func (m *mockResource) GetNamespace() string                     { return m.namespace }
func (m *mockResource) SetNamespace(namespace string)            { m.namespace = namespace }
func (m *mockResource) GetOwnerReference() models.OwnerReference { return m.ownerRef }
func (m *mockResource) SetOwnerReference(owner models.OwnerReference) {
	m.ownerRef = owner
}
func (m *mockResource) GetID() string                    { return m.kind + ":" + m.name + ":" + m.namespace }
func (m *mockResource) GetCreationTimestamp() time.Time  { return m.creationTimestamp }
func (m *mockResource) SetCreationTimestamp(t time.Time) { m.creationTimestamp = t }

func newRepoWithResume(t *testing.T) (*repository.InMemoryRepository, *mockResource) {
	t.Helper()
	repo := repository.NewInMemoryRepository()
	now := time.Now()
	ns := &mockResource{kind: "namespace", name: "default", creationTimestamp: now}
	if _, err := repo.Create(ns); err != nil {
		t.Fatalf("setup: create namespace: %v", err)
	}
	resume := &mockResource{kind: "resume", name: "test-resume", namespace: "default", creationTimestamp: now}
	if _, err := repo.Create(resume); err != nil {
		t.Fatalf("setup: create resume: %v", err)
	}
	return repo, resume
}

func TestInMemoryRepository_CreateAndDuplicate(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	now := time.Now()
	resource := &mockResource{kind: "resume", name: "r1", namespace: "default", creationTimestamp: now}

	msg, err := repo.Create(resource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "created") {
		t.Errorf("expected created message, got: %s", msg)
	}

	if _, err := repo.Create(resource); err == nil {
		t.Error("expected error creating duplicate resource")
	}
}

func TestInMemoryRepository_List(t *testing.T) {
	repo, resume := newRepoWithResume(t)

	resources, err := repo.List(resume.GetKind(), resume.GetName(), resume.GetNamespace())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 || resources[0].GetID() != resume.GetID() {
		t.Errorf("expected exactly resume in list, got %v", resources)
	}
}

func TestInMemoryRepository_ListInvalidKind(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	if _, err := repo.List("nonexistent", "", ""); err == nil {
		t.Error("expected error for invalid kind")
	}
}

func TestInMemoryRepository_ListEmptyNamespace(t *testing.T) {
	repo, _ := newRepoWithResume(t)
	resources, err := repo.List("resume", "", "invalid-ns")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(resources))
	}
}

func TestInMemoryRepository_Delete(t *testing.T) {
	repo, resume := newRepoWithResume(t)

	msg, err := repo.Delete(resume.GetKind(), resume.GetName(), resume.GetNamespace())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "deleted") {
		t.Errorf("expected deleted message, got: %s", msg)
	}
}

func TestInMemoryRepository_DeleteNonexistent(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	if _, err := repo.Delete("resume", "nonexistent", "default"); err == nil {
		t.Error("expected error deleting nonexistent resource")
	}
}

func TestInMemoryRepository_DeleteInvalidKind(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	if _, err := repo.Delete("invalid", "name", "default"); err == nil {
		t.Error("expected error for invalid kind")
	}
}

func TestInMemoryRepository_ListAcrossNamespaces(t *testing.T) {
	repo, _ := newRepoWithResume(t)
	now := time.Now()

	other := &mockResource{kind: "namespace", name: "other", creationTimestamp: now}
	if _, err := repo.Create(other); err != nil {
		t.Fatalf("setup: %v", err)
	}
	otherResume := &mockResource{kind: "resume", name: "r2", namespace: "other", creationTimestamp: now}
	if _, err := repo.Create(otherResume); err != nil {
		t.Fatalf("setup: %v", err)
	}

	resources, err := repo.List("namespace", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Errorf("expected 2 namespaces, got %d", len(resources))
	}
}

func TestInMemoryRepository_CreateStampsZeroTimestamp(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	r := &mockResource{kind: "resume", name: "no-ts", namespace: "default"}
	if _, err := repo.Create(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.GetCreationTimestamp().IsZero() {
		t.Error("Create did not auto-stamp a zero CreationTimestamp")
	}
}

func TestInMemoryRepository_CreatePreservesNonZeroTimestamp(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	fixed := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	r := &mockResource{kind: "resume", name: "with-ts", namespace: "default", creationTimestamp: fixed}
	if _, err := repo.Create(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stored, err := repo.List("resume", "with-ts", "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stored) != 1 || !stored[0].GetCreationTimestamp().Equal(fixed) {
		t.Errorf("Create overwrote CreationTimestamp: got %v, want %v",
			stored[0].GetCreationTimestamp(), fixed)
	}
}
