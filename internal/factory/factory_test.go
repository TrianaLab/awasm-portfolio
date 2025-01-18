package factory_test

import (
	"awasm-portfolio/internal/factory"
	"testing"
)

// Helper to create test data
func createTestData(name, namespace string) map[string]interface{} {
	return map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	}
}

// TestResourceFactory_Create tests the Create method for all supported resource kinds.
func TestResourceFactory_Create(t *testing.T) {
	factory := factory.NewResourceFactory()
	name := "testName"
	namespace := "testNamespace"

	// List of supported resource kinds
	kinds := []string{
		"profile",
		"namespace",
		"education",
		"experience",
		"contact",
		"certifications",
		"contributions",
		"skills",
	}

	// Test resource creation for each kind
	for _, kind := range kinds {
		resource := factory.Create(kind, createTestData(name, namespace))
		if resource == nil {
			t.Errorf("expected resource of kind %s to be created, but got nil", kind)
			continue
		}

		// Verify common properties
		if resource.GetKind() != kind {
			t.Errorf("expected kind %s, got %s", kind, resource.GetKind())
		}
		if resource.GetName() != name {
			t.Errorf("expected name %s, got %s", name, resource.GetName())
		}
		if kind != "namespace" && resource.GetNamespace() != namespace {
			t.Errorf("expected namespace %s, got %s", namespace, resource.GetNamespace())
		}
	}

	// Test unsupported resource kind
	unsupportedKind := "unsupported"
	resource := factory.Create(unsupportedKind, createTestData(name, namespace))
	if resource != nil {
		t.Errorf("expected nil for unsupported kind %s, but got resource", unsupportedKind)
	}
}

// TestResourceFactory_New verifies the creation of a ResourceFactory instance.
func TestResourceFactory_New(t *testing.T) {
	factory := factory.NewResourceFactory()
	if factory == nil {
		t.Fatal("expected non-nil factory instance, got nil")
	}
}
