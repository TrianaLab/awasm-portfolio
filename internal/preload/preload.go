package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
)

func PreloadData(repo *repository.InMemoryRepository) {
	// Preload namespaces
	repo.Create(&types.Namespace{Name: "default"})
	repo.Create(&types.Namespace{Name: "dev"})
	repo.Create(&types.Namespace{Name: "test"})

	// Preload profiles
	johnDoeProfile := &types.Profile{
		Name:      "john-doe",
		Namespace: "default",
	}
	repo.Create(johnDoeProfile)
	janeDoeProfile := &types.Profile{
		Name:      "jane-doe",
		Namespace: "dev",
	}
	repo.Create(janeDoeProfile)
	testUserProfile := &types.Profile{
		Name:      "test-user",
		Namespace: "test",
	}
	repo.Create(testUserProfile)

	// Preload certifications
	repo.Create(&types.Certifications{
		Name:      "john-doe-certifications",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      johnDoeProfile.GetKind(),
			Name:      johnDoeProfile.GetName(),
			Namespace: johnDoeProfile.GetNamespace(),
		},
		Certifications: []types.Certification{
			{
				Description: "AWS Certified Solutions Architect",
				Link:        "https://aws.amazon.com/certification/",
			},
		},
	})

	repo.Create(&types.Certifications{
		Name:      "jane-doe-certifications",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "jane-doe",
			Namespace: "dev",
		},
		Certifications: []types.Certification{
			{
				Description: "Google Cloud Certified Professional Cloud Architect",
				Link:        "https://cloud.google.com/certification/",
			},
		},
	})

	// Preload contacts
	johnDoeContact := &types.Contact{
		Name:      "john-doe-contact",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      johnDoeProfile.GetKind(),
			Name:      johnDoeProfile.GetName(),
			Namespace: johnDoeProfile.GetNamespace(),
		},
		Email:    "john.doe@example.com",
		LinkedIn: "https://linkedin.com/in/johndoe",
		GitHub:   "https://github.com/johndoe",
	}
	repo.Create(johnDoeContact)
	johnDoeProfile.Contact = *johnDoeContact // Link contact resource

	repo.Create(&types.Contact{
		Name:      "jane-doe-contact",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      janeDoeProfile.GetKind(),
			Name:      janeDoeProfile.GetName(),
			Namespace: janeDoeProfile.GetNamespace(),
		},
		Email:    "jane.doe@example.com",
		LinkedIn: "https://linkedin.com/in/janedoe",
		GitHub:   "https://github.com/janedoe",
	})

	// Preload contributions
	johnDoeContributions := &types.Contributions{
		Name:      "john-doe-contributions",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      johnDoeProfile.GetKind(),
			Name:      johnDoeProfile.GetName(),
			Namespace: johnDoeProfile.GetNamespace(),
		},
		Contributions: []types.Contribution{
			{
				Project:     "Open Source CLI Tool",
				Description: "Built a CLI tool for Kubernetes management.",
				Link:        "https://github.com/johndoe/cli-tool",
			},
		},
	}
	repo.Create(johnDoeContributions)
	johnDoeProfile.Contributions = *johnDoeContributions // Link contributions resource

	// Preload education
	janeDoeEducation := &types.Education{
		Name:      "jane-doe-education",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      janeDoeProfile.GetKind(),
			Name:      janeDoeProfile.GetName(),
			Namespace: janeDoeProfile.GetNamespace(),
		},
		Courses: []types.Course{
			{
				Title:       "Computer Science",
				Institution: "Stanford",
				Duration:    "2010-2014",
			},
		},
	}
	repo.Create(janeDoeEducation)
	janeDoeProfile.Education = *janeDoeEducation // Link education resource
}
