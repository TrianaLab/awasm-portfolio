package repository_test

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"fmt"
	"strings"
	"testing"
	"time"
)

// dummyResource is a simple implementation of models.Resource for testing.
type dummyResource struct {
	kind              string
	name              string
	namespace         string
	id                string
	creationTimestamp time.Time
	ownerReference    models.OwnerReference
}

func (d *dummyResource) GetKind() string                          { return d.kind }
func (d *dummyResource) GetName() string                          { return d.name }
func (d *dummyResource) SetName(name string)                      { d.name = name }
func (d *dummyResource) GetNamespace() string                     { return d.namespace }
func (d *dummyResource) SetNamespace(namespace string)            { d.namespace = namespace }
func (d *dummyResource) GetOwnerReference() models.OwnerReference { return d.ownerReference }
func (d *dummyResource) SetOwnerReference(owner models.OwnerReference) {
	d.ownerReference = owner
}
func (d *dummyResource) GetID() string {
	if d.id == "" {
		d.id = fmt.Sprintf("%s:%s:%s", d.kind, d.name, d.namespace)
	}
	return d.id
}
func (d *dummyResource) GetCreationTimestamp() time.Time  { return d.creationTimestamp }
func (d *dummyResource) SetCreationTimestamp(t time.Time) { d.creationTimestamp = t }

// Ensure dummyResource implements models.Resource.
var _ models.Resource = &dummyResource{}

func newResource(kind, name, namespace string) *dummyResource {
	return &dummyResource{kind: kind, name: name, namespace: namespace}
}

func TestCreateInvalidAndDuplicate(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	// Crear recurso válido
	validRes := newResource("resume", "user1", "ns1")
	msg, err := repo.Create(validRes)
	if err != nil {
		t.Fatalf("unexpected error creating valid resource: %v", err)
	}
	if !strings.Contains(msg, "resume/user1 created") {
		t.Errorf("unexpected create message: %s", msg)
	}

	// Crear recurso con tipo no soportado
	invalidRes := newResource("invalidKind", "bad", "ns1")
	msg, err = repo.Create(invalidRes)
	if err != nil {
		t.Fatalf("unexpected error creating resource with unsupported kind: %v", err)
	}
	if !strings.Contains(msg, "invalidKind/bad created") {
		t.Errorf("unexpected create message for invalid kind: %s", msg)
	}

	// Crear recurso duplicado
	_, err = repo.Create(validRes)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected duplicate creation error, got: %v", err)
	}
}

func TestListScenarios(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	// Poblar repositorio con un recurso válido
	validRes := newResource("resume", "user1", "ns1")
	_, _ = repo.Create(validRes)

	// Listar tipo de recurso no válido
	_, err := repo.List("unsupportedKind", "", "")
	if err == nil || !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error, got: %v", err)
	}

	// Listar recurso inexistente
	_, err = repo.List("resume", "nonexistent", "ns1")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error for unexisting resource, got: %v", err)
	}

	// Listar recurso válido existente
	results, err := repo.List("resume", "user1", "ns1")
	if err != nil {
		t.Fatalf("unexpected error listing existing resource: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for existing resource, got %d", len(results))
	}
}

func TestDeleteScenarios(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	// Populate repository with one valid resource.
	validRes := newResource("profile", "user1", "ns1")
	_, _ = repo.Create(validRes)

	// Delete invalid resource kind
	_, err := repo.Delete("unsupportedKind", "anything", "ns1")
	if err == nil || !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error on delete, got: %v", err)
	}

	// Delete unexisting resource
	_, err = repo.Delete("profile", "nonexistent", "ns1")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error on deleting nonexisting resource, got: %v", err)
	}

	// Delete existing resource
	delMsg, err := repo.Delete("profile", "user1", "ns1")
	if err != nil {
		t.Fatalf("unexpected error deleting existing resource: %v", err)
	}
	if !strings.Contains(delMsg, "profile/user1") {
		t.Errorf("delete message did not mention deleted resource: %s", delMsg)
	}
}

func TestCascadingNamespaceDelete(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	// Crear namespace y recursos dentro de él
	ns := newResource("namespace", "nsCascade", "")
	if _, err := repo.Create(ns); err != nil {
		t.Fatalf("failed to create namespace: %v", err)
	}

	work := newResource("work", "work1", "nsCascade")
	education := newResource("education", "edu1", "nsCascade")
	if _, err := repo.Create(work); err != nil {
		t.Fatalf("failed to create work: %v", err)
	}
	if _, err := repo.Create(education); err != nil {
		t.Fatalf("failed to create education: %v", err)
	}

	// Intentar borrado en cascada del namespace
	delMsg, err := repo.Delete("all", "", "nsCascade")
	if err != nil {
		t.Fatalf("unexpected error deleting existing resource: %v", err)
	}
	if !strings.Contains(delMsg, "namespace/nsCascade") {
		t.Errorf("delete message did not mention deleted resource: %s", delMsg)
	}

	// Confirmar que los recursos dentro del namespace fueron borrados
	_, err = repo.List("work", "work1", "nsCascade")
	if err == nil {
		t.Errorf("expected error listing deleted work resource, but resource still exists")
	}
	_, err = repo.List("education", "edu1", "nsCascade")
	if err == nil {
		t.Errorf("expected error listing deleted education resource, but resource still exists")
	}
}
