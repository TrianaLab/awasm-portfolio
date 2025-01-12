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
	johnDoeCertifications := &types.Certifications{
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
			{
				Description: "Certified Kubernetes Administrator",
				Link:        "https://www.cncf.io/certification/cka/",
			},
		},
	}
	repo.Create(johnDoeCertifications)
	johnDoeProfile.Certifications = *johnDoeCertifications

	janeDoeCertifications := &types.Certifications{
		Name:      "jane-doe-certifications",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      janeDoeProfile.GetKind(),
			Name:      janeDoeProfile.GetName(),
			Namespace: janeDoeProfile.GetNamespace(),
		},
		Certifications: []types.Certification{
			{
				Description: "Google Cloud Professional Cloud Architect",
				Link:        "https://cloud.google.com/certification/",
			},
		},
	}
	repo.Create(janeDoeCertifications)
	janeDoeProfile.Certifications = *janeDoeCertifications

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
	johnDoeProfile.Contact = *johnDoeContact

	janeDoeContact := &types.Contact{
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
	}
	repo.Create(janeDoeContact)
	janeDoeProfile.Contact = *janeDoeContact

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
			{
				Project:     "Dashboard for Kubernetes",
				Description: "Developed a web dashboard for visualizing Kubernetes clusters.",
				Link:        "https://github.com/johndoe/k8s-dashboard",
			},
		},
	}
	repo.Create(johnDoeContributions)
	johnDoeProfile.Contributions = *johnDoeContributions

	janeDoeContributions := &types.Contributions{
		Name:      "jane-doe-contributions",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      janeDoeProfile.GetKind(),
			Name:      janeDoeProfile.GetName(),
			Namespace: janeDoeProfile.GetNamespace(),
		},
		Contributions: []types.Contribution{
			{
				Project:     "Cloud Monitoring Tool",
				Description: "Created a monitoring tool for cloud infrastructure.",
				Link:        "https://github.com/janedoe/cloud-monitor",
			},
		},
	}
	repo.Create(janeDoeContributions)
	janeDoeProfile.Contributions = *janeDoeContributions

	// Preload education
	johnDoeEducation := &types.Education{
		Name:      "john-doe-education",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      johnDoeProfile.GetKind(),
			Name:      johnDoeProfile.GetName(),
			Namespace: johnDoeProfile.GetNamespace(),
		},
		Courses: []types.Course{
			{
				Title:       "Computer Science",
				Institution: "MIT",
				Duration:    "2010-2014",
			},
			{
				Title:       "Data Science",
				Institution: "Harvard",
				Duration:    "2015-2016",
			},
		},
	}
	repo.Create(johnDoeEducation)
	johnDoeProfile.Education = *johnDoeEducation

	// Preload experiences
	johnDoeExperience := &types.Experience{
		Name:      "john-doe-experience",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      johnDoeProfile.GetKind(),
			Name:      johnDoeProfile.GetName(),
			Namespace: johnDoeProfile.GetNamespace(),
		},
		Jobs: []types.Job{
			{
				Title:       "Software Engineer",
				Company:     "TechCorp",
				Description: "Worked on cloud infrastructure tools.",
				Duration:    "2015-2020",
			},
			{
				Title:       "Lead Developer",
				Company:     "Innovatech",
				Description: "Led a team building AI-driven solutions.",
				Duration:    "2020-Present",
			},
		},
	}
	repo.Create(johnDoeExperience)
	johnDoeProfile.Experience = *johnDoeExperience

	// Save profiles with linked resources
	repo.Create(johnDoeProfile)
	repo.Create(janeDoeProfile)
	repo.Create(testUserProfile)
}
