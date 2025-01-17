package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
)

var resources = []models.Resource{
	// Namespaces
	&types.Namespace{Name: "default"},
	&types.Namespace{Name: "dev"},
	&types.Namespace{Name: "test"},

	// Profiles with child resources
	&types.Profile{
		Name:      "john-doe",
		Namespace: "default",
		Certifications: types.Certifications{
			Name:      "john-doe-certifications",
			Namespace: "default",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "john-doe",
				Namespace: "default",
			},
			Certifications: []types.Certification{
				{Description: "AWS Certified Solutions Architect", Link: "https://aws.amazon.com/certification/"},
				{Description: "Certified Kubernetes Administrator", Link: "https://www.cncf.io/certification/cka/"},
			},
		},
		Contact: types.Contact{
			Name:      "john-doe-contact",
			Namespace: "default",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "john-doe",
				Namespace: "default",
			},
			Email:    "john.doe@example.com",
			Linkedin: "https://linkedin.com/in/johndoe",
			Github:   "https://github.com/johndoe",
		},
		Contributions: types.Contributions{
			Name:      "john-doe-contributions",
			Namespace: "default",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "john-doe",
				Namespace: "default",
			},
			Contributions: []types.Contribution{
				{Project: "Open Source CLI Tool", Description: "Built a CLI tool for Kubernetes management.", Link: "https://github.com/johndoe/cli-tool"},
				{Project: "Dashboard for Kubernetes", Description: "Developed a web dashboard for visualizing Kubernetes clusters.", Link: "https://github.com/johndoe/k8s-dashboard"},
			},
		},
		Education: types.Education{
			Name:      "john-doe-education",
			Namespace: "default",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "john-doe",
				Namespace: "default",
			},
			Courses: []types.Course{
				{Title: "Computer Science", Institution: "MIT", Duration: "2010-2014"},
				{Title: "Data Science", Institution: "Harvard", Duration: "2015-2016"},
			},
		},
		Skills: types.Skills{
			Name:      "john-doe-skills",
			Namespace: "default",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "john-doe",
				Namespace: "default",
			},
			Skills: []types.Skill{
				{Competence: "Go", Proficiency: "Expert"},
				{Competence: "Kubernetes", Proficiency: "Advanced"},
				{Competence: "Docker", Proficiency: "Advanced"},
				{Competence: "Cloud Infrastructure", Proficiency: "Expert"},
				{Competence: "Web Development", Proficiency: "Intermediate"},
			},
		},
		Experience: types.Experience{
			Name:      "john-doe-experience",
			Namespace: "default",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "john-doe",
				Namespace: "default",
			},
			Jobs: []types.Job{
				{Title: "Software Engineer", Company: "TechCorp", Description: "Worked on cloud infrastructure tools.", Duration: "2015-2020"},
				{Title: "Lead Developer", Company: "Innovatech", Description: "Led a team building AI-driven solutions.", Duration: "2020-Present"},
			},
		},
	},
	&types.Profile{
		Name:      "jane-doe",
		Namespace: "dev",
		Certifications: types.Certifications{
			Name:      "jane-doe-certifications",
			Namespace: "dev",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "jane-doe",
				Namespace: "dev",
			},
			Certifications: []types.Certification{
				{Description: "Google Cloud Professional Cloud Architect", Link: "https://cloud.google.com/certification/"},
			},
		},
		Contact: types.Contact{
			Name:      "jane-doe-contact",
			Namespace: "dev",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "jane-doe",
				Namespace: "dev",
			},
			Email:    "jane.doe@example.com",
			Linkedin: "https://linkedin.com/in/janedoe",
			Github:   "https://github.com/janedoe",
		},
		Contributions: types.Contributions{
			Name:      "jane-doe-contributions",
			Namespace: "dev",
			OwnerRef: models.OwnerReference{
				Kind:      "Profile",
				Name:      "jane-doe",
				Namespace: "dev",
			},
			Contributions: []types.Contribution{
				{Project: "Cloud Monitoring Tool", Description: "Created a monitoring tool for cloud infrastructure.", Link: "https://github.com/janedoe/cloud-monitor"},
			},
		},
	},
	&types.Profile{
		Name:      "test-user",
		Namespace: "test",
	},
}
