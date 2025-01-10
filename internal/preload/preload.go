package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
	"fmt"
)

// PreloadData preloads namespaces, profiles, and associated resources
func PreloadData(repo *repository.InMemoryRepository) {
	// Preload namespaces
	namespaces := []string{"default", "dev", "test"}
	for _, ns := range namespaces {
		repo.Create("namespace", &types.Namespace{
			Name: ns,
		})
		fmt.Printf("Preloaded namespace: %s\n", ns)
	}

	// Preload profiles in different namespaces
	profiles := []struct {
		Name      string
		Namespace string
	}{
		{"john-doe", "default"},
		{"jane-doe", "default"},
		{"dev-user", "dev"},
		{"test-user", "test"},
	}

	for _, profile := range profiles {
		repo.Create("profile", &types.Profile{
			Name:      profile.Name,
			Namespace: profile.Namespace,
			OwnerRefs: []models.OwnerReference{
				{Kind: "namespace", Name: profile.Namespace}, // Namespace as owner
			},
		})
		fmt.Printf("Preloaded profile: %s in namespace: %s\n", profile.Name, profile.Namespace)
	}

	// Preload resources with profiles as their owners
	repo.Create("experience", &types.Experience{
		Name:      "senior-developer",
		Namespace: "default",
		OwnerRefs: []models.OwnerReference{
			{Kind: "profile", Name: "john-doe"}, // Profile as owner
		},
		Jobs: []types.Job{
			{Title: "Senior Developer", Company: "TechCorp", Duration: "3 years", Description: "Backend development"},
		},
	})
	fmt.Printf("Preloaded experience: senior-developer owned by profile: john-doe\n")

	repo.Create("education", &types.Education{
		Name:      "bachelor-cs",
		Namespace: "default",
		OwnerRefs: []models.OwnerReference{
			{Kind: "profile", Name: "jane-doe"}, // Profile as owner
		},
		Courses: []types.Course{
			{Title: "Computer Science", Institution: "XYZ University", Duration: "4 years"},
		},
	})
	fmt.Printf("Preloaded education: bachelor-cs owned by profile: jane-doe\n")

	repo.Create("certifications", &types.Certifications{
		Name:      "aws-certified",
		Namespace: "dev",
		OwnerRefs: []models.OwnerReference{
			{Kind: "profile", Name: "dev-user"}, // Profile as owner
		},
		Certifications: []types.Certification{
			{Description: "AWS Certified Solutions Architect", Link: "https://aws.amazon.com"},
		},
	})
	fmt.Printf("Preloaded certification: aws-certified owned by profile: dev-user\n")

	repo.Create("contributions", &types.Contributions{
		Name:      "open-source-projects",
		Namespace: "test",
		OwnerRefs: []models.OwnerReference{
			{Kind: "profile", Name: "test-user"}, // Profile as owner
		},
		Contributions: []types.Contribution{
			{Project: "CLI Tool", Description: "A Kubernetes-like CLI for managing resources", Link: "https://github.com/example/cli-tool"},
		},
	})
	fmt.Printf("Preloaded contributions: open-source-projects owned by profile: test-user\n")

	repo.Create("skills", &types.Skills{
		Name:      "programming-languages",
		Namespace: "default",
		OwnerRefs: []models.OwnerReference{
			{Kind: "profile", Name: "john-doe"}, // Profile as owner
		},
		Skills: []types.Skill{
			{Name: "Go", Proficiency: "Expert"},
			{Name: "JavaScript", Proficiency: "Intermediate"},
		},
	})
	fmt.Printf("Preloaded skills: programming-languages owned by profile: john-doe\n")

	repo.Create("skills", &types.Skills{
		Name:      "devops-tools",
		Namespace: "dev",
		OwnerRefs: []models.OwnerReference{
			{Kind: "profile", Name: "dev-user"}, // Profile as owner
		},
		Skills: []types.Skill{
			{Name: "Kubernetes", Proficiency: "Advanced"},
			{Name: "Docker", Proficiency: "Expert"},
		},
	})
	fmt.Printf("Preloaded skills: devops-tools owned by profile: dev-user\n")
}
