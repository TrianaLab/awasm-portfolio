package preload

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
	"fmt"
)

func PreloadData(repo *repository.InMemoryRepository) {
	factory := factory.NewResourceFactory()

	// Preload namespaces
	namespaces := []string{"default", "dev", "test"}
	for _, ns := range namespaces {
		namespace := factory.Create("namespace", map[string]interface{}{
			"name": ns,
		})
		repo.Create(namespace)
		fmt.Printf("Preloaded namespace: %s\n", ns)
	}

	// Preload profiles with associated data
	profiles := []struct {
		Name           string
		Namespace      string
		Certifications []types.Certification
		Contact        types.Contact
		Experience     []types.Job
		Education      []types.Course
		Skills         []types.Skill
	}{
		{
			Name:      "john-doe",
			Namespace: "default",
			Certifications: []types.Certification{
				{Description: "AWS Certified Solutions Architect", Link: "https://aws.amazon.com"},
				{Description: "Certified Kubernetes Administrator", Link: "https://kubernetes.io"},
			},
			Contact: types.Contact{
				Email:    "john.doe@example.com",
				LinkedIn: "linkedin.com/in/john-doe",
				GitHub:   "github.com/johndoe",
			},
			Experience: []types.Job{
				{Title: "Backend Developer", Company: "TechCorp", Duration: "2 years", Description: "Developed microservices in Go."},
				{Title: "DevOps Engineer", Company: "CloudOps", Duration: "3 years", Description: "Implemented CI/CD pipelines."},
			},
			Education: []types.Course{
				{Title: "BSc in Computer Science", Institution: "XYZ University", Duration: "4 years"},
			},
			Skills: []types.Skill{
				{Name: "Go", Proficiency: "Expert"},
				{Name: "Docker", Proficiency: "Advanced"},
			},
		},
		{
			Name:      "jane-doe",
			Namespace: "dev",
			Certifications: []types.Certification{
				{Description: "Google Cloud Professional Data Engineer", Link: "https://cloud.google.com"},
			},
			Contact: types.Contact{
				Email:    "jane.doe@example.com",
				LinkedIn: "linkedin.com/in/jane-doe",
				GitHub:   "github.com/janedoe",
			},
			Experience: []types.Job{
				{Title: "Frontend Developer", Company: "WebStudio", Duration: "1 year", Description: "Built responsive web apps."},
			},
			Education: []types.Course{
				{Title: "MSc in Data Science", Institution: "ABC University", Duration: "2 years"},
			},
			Skills: []types.Skill{
				{Name: "React", Proficiency: "Advanced"},
				{Name: "Python", Proficiency: "Intermediate"},
			},
		},
		{
			Name:      "dev-user",
			Namespace: "test",
			Certifications: []types.Certification{
				{Description: "Red Hat Certified Engineer", Link: "https://redhat.com"},
			},
			Contact: types.Contact{
				Email:    "dev.user@example.com",
				LinkedIn: "linkedin.com/in/dev-user",
				GitHub:   "github.com/devuser",
			},
			Experience: []types.Job{
				{Title: "Full Stack Developer", Company: "CodeBase", Duration: "5 years", Description: "Developed full-stack applications."},
			},
			Education: []types.Course{
				{Title: "Diploma in Software Engineering", Institution: "PQR Institute", Duration: "3 years"},
			},
			Skills: []types.Skill{
				{Name: "JavaScript", Proficiency: "Advanced"},
				{Name: "Kubernetes", Proficiency: "Expert"},
			},
		},
	}

	for _, profile := range profiles {
		resource := factory.Create("profile", map[string]interface{}{
			"name":      profile.Name,
			"namespace": profile.Namespace,
		})
		repo.Create(resource)

		certifications := factory.Create("certifications", map[string]interface{}{
			"name":      profile.Name + "-certifications",
			"namespace": profile.Namespace,
		})
		repo.Create(certifications)

		contact := factory.Create("contact", map[string]interface{}{
			"name":      profile.Name + "-contact",
			"namespace": profile.Namespace,
		})
		repo.Create(contact)

		fmt.Printf("Preloaded profile: %s in namespace: %s\n", profile.Name, profile.Namespace)
	}
}
