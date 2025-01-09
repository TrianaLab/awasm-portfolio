package preload

import (
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
)

func PreloadData(repo *repository.InMemoryRepository) {
	// Preload profiles
	repo.Create("profile", &types.Profile{
		Name:      "john-doe",
		Namespace: "default",
		OwnerRefs: nil,
	})

	// Preload namespaces
	repo.Create("namespace", &types.Namespace{
		Name: "default",
	})

	// Preload education
	repo.Create("education", &types.Education{
		Name:      "bachelor-cs",
		Namespace: "default",
		Courses: []types.Course{
			{Title: "Computer Science", Institution: "XYZ University", Duration: "4 years"},
		},
	})

	// Preload experience
	repo.Create("experience", &types.Experience{
		Name:      "senior-developer",
		Namespace: "default",
		Jobs: []types.Job{
			{Title: "Senior Developer", Company: "TechCorp", Description: "Backend development", Duration: "3 years"},
		},
	})

	// Preload skills
	repo.Create("skills", &types.Skills{
		Name:      "programming-languages",
		Namespace: "default",
		Skills: []types.Skill{
			{Name: "Go", Proficiency: "Expert"},
			{Name: "JavaScript", Proficiency: "Intermediate"},
		},
	})

	// Preload certifications
	repo.Create("certifications", &types.Certifications{
		Name:      "aws-certified",
		Namespace: "default",
		Certifications: []types.Certification{
			{Description: "AWS Certified Solutions Architect", Link: "https://aws.amazon.com"},
		},
	})

	// Preload contributions
	repo.Create("contributions", &types.Contributions{
		Name:      "open-source-projects",
		Namespace: "default",
		Contributions: []types.Contribution{
			{Project: "CLI Tool", Description: "A Kubernetes-like CLI for managing resources", Link: "https://github.com/example/cli-tool"},
		},
	})

	// Preload contact information
	repo.Create("contact", &types.Contact{
		Name:      "john-contact",
		Namespace: "default",
		Email:     "john.doe@example.com",
		LinkedIn:  "https://linkedin.com/in/johndoe",
		GitHub:    "https://github.com/johndoe",
	})
}
