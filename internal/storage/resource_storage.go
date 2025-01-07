package storage

import "awasm-portfolio/internal/models"

func PreloadData(resourceManager *ResourceManager) {
	// Create a namespaced Contact resource
	contact := models.ResourceBase{
		Name:      "contact1",
		Namespace: "test",
		Data: map[string]interface{}{
			"email":    "test.user@example.com",
			"linkedin": "linkedin.com/in/testuser",
			"github":   "github.com/testuser",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("contact", contact)

	// Create namespaced Skill resources
	skills := []models.ResourceBase{
		{
			Name:       "C",
			Namespace:  "test",
			Data:       nil, // No additional data for skill resources
			Namespaced: true,
		},
		{
			Name:       "Java",
			Namespace:  "test",
			Data:       nil,
			Namespaced: true,
		},
		{
			Name:       "Go",
			Namespace:  "test",
			Data:       nil,
			Namespaced: true,
		},
	}

	// Add skills to the resource manager
	for _, skill := range skills {
		_ = resourceManager.Create("skill", skill)
	}

	// Include the Contact and Skills as part of the Profile's Data
	profile := models.ResourceBase{
		Name:      "test",
		Namespace: "",
		Data: map[string]interface{}{
			"firstname": "Test",
			"lastname":  "User",
			"contact":   contact, // Embed the Contact ResourceBase
			"skills":    skills,  // Embed the list of Skill ResourceBases
		},
		Namespaced: false,
	}
	_ = resourceManager.Create("profile", profile)
}
