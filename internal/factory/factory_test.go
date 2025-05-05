package factory_test

import (
	"awasm-portfolio/internal/factory"
	"testing"
)

func createTestData(name, namespace string) map[string]interface{} {
	return map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	}
}

func TestResourceFactory_Create(t *testing.T) {
	factory := factory.NewResourceFactory()
	name := "testName"
	namespace := "testNamespace"

	// Lista actualizada de tipos de recursos soportados
	kinds := []string{
		"resume",
		"work",
		"volunteer",
		"education",
		"award",
		"certificate",
		"publication",
		"skill",
		"language",
		"interest",
		"reference",
		"project",
	}

	for _, kind := range kinds {
		t.Run(kind, func(t *testing.T) {
			resource := factory.Create(kind, createTestData(name, namespace))
			if resource == nil {
				t.Fatalf("expected resource of kind %s to be created, but got nil", kind)
			}

			// Verificar propiedades comunes
			if resource.GetKind() != kind {
				t.Errorf("expected kind %s, got %s", kind, resource.GetKind())
			}

			// El nombre puede variar según el tipo debido a que usamos faker
			if resource.GetName() == "" {
				t.Error("expected non-empty name")
			}

			if resource.GetNamespace() != namespace {
				t.Errorf("expected namespace %s, got %s", namespace, resource.GetNamespace())
			}

			// Verificar que el timestamp de creación no esté vacío
			if resource.GetCreationTimestamp().IsZero() {
				t.Error("expected non-zero creation timestamp")
			}

			// Verificar que la referencia del propietario esté establecida
			if resource.GetOwnerReference().Name != namespace {
				t.Errorf("expected owner reference name %s, got %s", namespace, resource.GetOwnerReference().Name)
			}
		})
	}

	t.Run("unsupported", func(t *testing.T) {
		unsupportedKind := "unsupported"
		resource := factory.Create(unsupportedKind, createTestData(name, namespace))
		if resource != nil {
			t.Errorf("expected nil for unsupported kind %s, but got resource", unsupportedKind)
		}
	})
}

func TestResourceFactory_New(t *testing.T) {
	factory := factory.NewResourceFactory()
	if factory == nil {
		t.Fatal("expected non-nil factory instance, got nil")
	}
}
