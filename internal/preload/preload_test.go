package preload_test

import (
	"awasm-portfolio/internal/preload"
	"awasm-portfolio/internal/repository"
	"testing"
)

// TestPreloadData ensures that PreloadData correctly preloads resources into the repository.
func TestPreloadData(t *testing.T) {
	// Create a new in-memory repository
	repo := repository.NewInMemoryRepository()

	// Preload the data
	preload.PreloadData(repo)

	// Define the expected resources to verify
	expectedKinds := []string{
		"namespace",
		"award",
		"basics",
		"certificate",
		"education",
		"interest",
		"language",
		"resume",
		"skill",
		"volunteer",
		"work",
	}

	// Verify that each expected resource kind exists in the repository
	for _, kind := range expectedKinds {
		resources, err := repo.List(kind, "", "")
		if err != nil {
			t.Fatalf("unexpected error listing resources of kind %s: %v", kind, err)
		}

		if len(resources) == 0 {
			t.Errorf("expected resources of kind %s to be preloaded, but none found", kind)
		}
	}

	// Verify a specific resource: namespace
	namespace, err := repo.List("namespace", "default", "")
	if err != nil {
		t.Fatalf("unexpected error listing namespace resources: %v", err)
	}
	if len(namespace) != 1 {
		t.Fatalf("expected exactly 1 namespace resource, found %d", len(namespace))
	}
	if namespace[0].GetName() != "default" {
		t.Errorf("expected namespace resource name to be 'default', got %s", namespace[0].GetName())
	}
}
