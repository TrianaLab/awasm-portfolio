package storage

import "awasm-portfolio/internal/models"

func PreloadData(resourceManager *ResourceManager) {
	// Preload Namespace
	namespace := models.ResourceBase{
		Name:       "default",
		Namespace:  "",
		Namespaced: false,
	}
	_ = resourceManager.Create("namespace", namespace)

	// Preload Profiles
	profile := models.ResourceBase{
		Name:      "john_doe",
		Namespace: "default",
		Data: map[string]interface{}{
			"firstname":     "John",
			"lastname":      "Doe",
			"description":   "Experienced Software Engineer",
			"role":          "Senior Developer",
			"experience":    []models.ResourceBase{}, // Placeholder for linked experiences
			"contribution":  []models.ResourceBase{}, // Placeholder for linked contributions
			"education":     []models.ResourceBase{}, // Placeholder for linked educations
			"certification": []models.ResourceBase{}, // Placeholder for linked certifications
			"project":       []models.ResourceBase{}, // Placeholder for linked projects
			"skill":         []models.ResourceBase{}, // Placeholder for linked skills
			"contact": models.ResourceBase{
				Name: "contact1",
			}, // Linked contact
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("profile", profile)

	// Preload Experiences
	experience := models.ResourceBase{
		Name:      "experience1",
		Namespace: "default",
		Data: map[string]interface{}{
			"role":       "Software Engineer",
			"company":    "TechCorp",
			"start_date": "2020-01-01",
			"end_date":   "2023-01-01",
			"skill":      []models.ResourceBase{}, // Placeholder for linked skills
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("experience", experience)

	// Preload Contributions
	contribution := models.ResourceBase{
		Name:      "contribution1",
		Namespace: "default",
		Data: map[string]interface{}{
			"project":     "Open Source Library",
			"description": "Contributed to an open-source project",
			"link":        "https://github.com/example/project",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("contribution", contribution)

	// Preload Educations
	education := models.ResourceBase{
		Name:      "education1",
		Namespace: "default",
		Data: map[string]interface{}{
			"institution": "University of Tech",
			"degree":      "Bachelor's",
			"title":       "Computer Science",
			"start_date":  "2015-09-01",
			"end_date":    "2019-06-01",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("education", education)

	// Preload Certifications
	certification := models.ResourceBase{
		Name:      "certification1",
		Namespace: "default",
		Data: map[string]interface{}{
			"description": "AWS Certified Solutions Architect",
			"link":        "https://aws.amazon.com/certification/",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("certification", certification)

	// Preload Projects
	project := models.ResourceBase{
		Name:      "project1",
		Namespace: "default",
		Data: map[string]interface{}{
			"description": "Built a scalable web application",
			"link":        "https://github.com/example/scalable-webapp",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("project", project)

	// Preload Skills
	skill := models.ResourceBase{
		Name:       "go_programming",
		Namespace:  "default",
		Data:       nil, // No additional data
		Namespaced: true,
	}
	_ = resourceManager.Create("skill", skill)

	// Preload Contacts
	contact := models.ResourceBase{
		Name:      "contact1",
		Namespace: "default",
		Data: map[string]interface{}{
			"email":    "john.doe@example.com",
			"github":   "https://github.com/johndoe",
			"linkedin": "https://linkedin.com/in/johndoe",
		},
		Namespaced: true,
	}
	_ = resourceManager.Create("contact", contact)
}
