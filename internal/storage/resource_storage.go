package storage

import "awasm-portfolio/internal/models"

func PreloadData(resourceManager *ResourceManager) {
	// Preload a Namespace resource
	namespace := models.ResourceBase{
		Name:       "default",
		Namespace:  "", // Cluster-wide
		Namespaced: false,
	}
	_ = resourceManager.Create("namespace", namespace)

	// Adjusted Profile to be namespaced
	profile := models.ResourceBase{
		Name:      "test",
		Namespace: "default", // Now namespaced
		Data: map[string]interface{}{
			"firstname": "Test",
			"lastname":  "User",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("profile", profile)
}
