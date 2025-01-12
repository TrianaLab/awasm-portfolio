package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
)

func PreloadData(repo *repository.InMemoryRepository) {
	// Preload namespaces
	namespaces := []types.Namespace{
		{Name: "default"},
		{Name: "dev"},
		{Name: "test"},
	}
	for _, ns := range namespaces {
		_, _ = repo.Create(&ns)
	}

	// Preload a profile and child resources
	profile := &types.Profile{
		Name:      "john-doe",
		Namespace: "default",
		OwnerRef:  models.OwnerReference{Kind: "namespace", Name: "default"},
		Certifications: types.Certifications{
			Name:      "john-doe-certifications",
			Namespace: "default",
			OwnerRef:  models.OwnerReference{Kind: "profile", Name: "john-doe", Namespace: "default"},
			Certifications: []types.Certification{
				{Description: "AWS Certified Solutions Architect", Link: "https://aws.amazon.com/certification/"},
			},
		},
		Contact: types.Contact{
			Name:      "john-doe-contact",
			Namespace: "default",
			OwnerRef:  models.OwnerReference{Kind: "profile", Name: "john-doe", Namespace: "default"},
			Email:     "john.doe@example.com",
			LinkedIn:  "https://linkedin.com/in/johndoe",
			GitHub:    "https://github.com/johndoe",
		},
		Contributions: types.Contributions{
			Name:      "john-doe-contributions",
			Namespace: "default",
			OwnerRef:  models.OwnerReference{Kind: "profile", Name: "john-doe", Namespace: "default"},
			Contributions: []types.Contribution{
				{Project: "Open Source CLI Tool", Description: "Built a CLI tool for Kubernetes management.", Link: "https://github.com/johndoe/cli-tool"},
			},
		},
		Education: types.Education{
			Name:      "john-doe-education",
			Namespace: "default",
			OwnerRef:  models.OwnerReference{Kind: "profile", Name: "john-doe", Namespace: "default"},
			Courses: []types.Course{
				{Title: "Computer Science", Institution: "MIT", Duration: "2015-2019"},
			},
		},
		Experience: types.Experience{
			Name:      "john-doe-experience",
			Namespace: "default",
			OwnerRef:  models.OwnerReference{Kind: "profile", Name: "john-doe", Namespace: "default"},
			Jobs: []types.Job{
				{
					Title:       "Senior Software Engineer",
					Company:     "TechCorp",
					Duration:    "2020-2023",
					Description: "Led the development of scalable cloud applications.",
				},
			},
		},
		Skills: types.Skills{
			Name:      "john-doe-skills",
			Namespace: "default",
			OwnerRef:  models.OwnerReference{Kind: "profile", Name: "john-doe", Namespace: "default"},
			Skills: []types.Skill{
				{Name: "Go", Proficiency: "Expert"},
				{Name: "Kubernetes", Proficiency: "Advanced"},
			},
		},
	}

	// Create the profile and child resources
	_, _ = repo.Create(profile)
	_, _ = repo.Create(&profile.Certifications)
	_, _ = repo.Create(&profile.Contact)
	_, _ = repo.Create(&profile.Contributions)
	_, _ = repo.Create(&profile.Education)
	_, _ = repo.Create(&profile.Experience)
	_, _ = repo.Create(&profile.Skills)
}
