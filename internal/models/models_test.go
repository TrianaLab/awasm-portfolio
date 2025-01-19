package models_test

import (
	"awasm-portfolio/internal/models"
	"strings"
	"testing"
	"time"
)         string
	ownerReference    models.OwnerReference
	creationTimestamp time.Time
}

func (d *DummyResource) GetKind() string                               { return d.kind }
func (d *DummyResource) GetName() string                               { return d.name }
func (d *DummyResource) SetName(name string)                           { d.name = name }
func (d *DummyResource) GetNamespace() string                          { return d.namespace }
func (d *DummyResource) SetNamespace(namespace string)                 { d.namespace = namespace }
func (d *DummyResource) GetOwnerReference() models.OwnerReference      { return d.ownerReference }
func (d *DummyResource) SetOwnerReference(owner models.OwnerReference) { d.ownerReference = owner }
func (d *DummyResource) GetID() string {
	return strings.ToLower(d.kind + ":" + d.namespace + ":" + d.name)
}
func (d *DummyResource) GetCreationTimestamp() time.Time  { return d.creationTimestamp }
func (d *DummyResource) SetCreationTimestamp(t time.Time) { d.creationTimestamp = t }

// TestOwnerReference_GetID tests the GetID method of OwnerReference.
func TestOwnerReference_GetID(t *testing.T) {
	ref := models.OwnerReference{
		Kind:      "TestKind",
		Name:      "TestName",
		Namespace: "TestNamespace",
	}

	expectedID := "testkind:testname:testnamespace"
	actualID := ref.GetID()
	if actualID != expectedID {
		t.Errorf("expected ID %q, got %q", expectedID, actualID)
	}
}

// TestOwnerReference_GetName tests the GetName method of OwnerReference.
func TestOwnerReference_GetName(t *testing.T) {
	ref := models.OwnerReference{
		Kind:      "TestKind",
		Name:      "TestName",
		Namespace: "TestNamespace",
	}

	// Without Owner
	if ref.GetName() != "TestName" {
		t.Errorf("expected name %q, got %q", "TestName", ref.GetName())
	}

	// With Owner
	owner := &DummyResource{name: "OwnerName"}
	ref.Owner = owner
	if ref.GetName() != "OwnerName" {
		t.Errorf("expected owner name %q, got %q", "OwnerName", ref.GetName())
	}
}

// TestDummyResource tests the implementation of the Resource interface.
func TestDummyResource(t *testing.T) {
	resource := &DummyResource{
		kind:      "TestKind",
		name:      "TestName",
		namespace: "TestNamespace",
	}

	// Test GetKind
	if resource.GetKind() != "TestKind" {
		t.Errorf("expected kind %q, got %q", "TestKind", resource.GetKind())
	}

	// Test GetName and SetName
	resource.SetName("NewName")
	if resource.GetName() != "NewName" {
		t.Errorf("expected name %q, got %q", "NewName", resource.GetName())
	}

	// Test GetNamespace and SetNamespace
	resource.SetNamespace("NewNamespace")
	if resource.GetNamespace() != "NewNamespace" {
		t.Errorf("expected namespace %q, got %q", "NewNamespace", resource.GetNamespace())
	}

	// Test GetOwnerReference and SetOwnerReference
	ownerRef := models.OwnerReference{
		Kind:      "OwnerKind",
		Name:      "OwnerName",
		Namespace: "OwnerNamespace",
	}
	resource.SetOwnerReference(ownerRef)
	if resource.GetOwnerReference().GetID() != ownerRef.GetID() {
		t.Errorf("expected owner reference ID %q, got %q", ownerRef.GetID(), resource.GetOwnerReference().GetID())
	}

	// Test GetID
	expectedID := "testkind:newnamespace:newname"
	if resource.GetID() != expectedID {
		t.Errorf("expected ID %q, got %q", expectedID, resource.GetID())
	}

	// Test GetCreationTimestamp and SetCreationTimestamp
	timestamp := time.Now()
	resource.SetCreationTimestamp(timestamp)
	if !resource.GetCreationTimestamp().Equal(timestamp) {
		t.Errorf("expected timestamp %v, got %v", timestamp, resource.GetCreationTimestamp())
	}
}
