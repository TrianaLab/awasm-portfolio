package types

import (
	"awasm-portfolio/internal/models"
	"testing"
	"time"
)

func TestResumeWithAllFields(t *testing.T) {

	timestamp := time.Now()
	ownerRef := models.OwnerReference{
		Kind: "Namespace",
		Name: "default",
	}

	resume := &Resume{
		Kind:              "resume",
		Name:              "JohnDoeResume",
		Namespace:         "default",
		OwnerRef:          ownerRef,
		CreationTimestamp: timestamp,
		Basics: Basics{
			Name:    "John Doe",
			Label:   "Programmer",
			Image:   "https://example.com/image.jpg",
			Email:   "john.doe@example.com",
			Phone:   "123-456-7890",
			Url:     "https://johndoe.com",
			Summary: "A summary of John Doe",
			Location: Location{
				Address:     "123 Main St",
				PostalCode:  "12345",
				City:        "San Francisco",
				CountryCode: "US",
				Region:      "California",
			},
			Profiles: []Profile{
				{
					Network:  "Twitter",
					Username: "johndoe",
					Url:      "https://twitter.com/johndoe",
				},
			},
		},
		Work: []Work{
			{
				Kind:              "work",
				Name:              "Company",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Position:          "President",
				URL:               "https://company.com",
				StartDate:         "2013-01-01",
				EndDate:           "2014-01-01",
				Summary:           "Description of work",
				Highlights:        []string{"Started the company"},
			},
		},
		Volunteer: []Volunteer{
			{
				Kind:              "volunteer",
				Name:              "Organization",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Position:          "Volunteer",
				URL:               "https://organization.com",
				StartDate:         "2012-01-01",
				EndDate:           "2013-01-01",
				Summary:           "Description of volunteer work",
				Highlights:        []string{"Awarded 'Volunteer of the Month'"},
			},
		},
		Education: []Education{
			{
				Kind:              "education",
				Name:              "University",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				URL:               "https://university.com",
				Area:              "Software Development",
				StudyType:         "Bachelor",
				StartDate:         "2011-01-01",
				EndDate:           "2013-01-01",
				Score:             "4.0",
				Courses:           []string{"DB1101 - Basic SQL"},
			},
		},
		Awards: []Award{
			{
				Kind:              "award",
				Name:              "Best Developer",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Date:              "2014-11-01",
				Awarder:           "Tech Company",
				Summary:           "Awarded for outstanding performance",
			},
		},
		Certificates: []Certificate{
			{
				Kind:              "certificate",
				Name:              "Certified Kubernetes Administrator",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Date:              "2021-11-07",
				Issuer:            "Cloud Native Computing Foundation",
				URL:               "https://certificate.com",
			},
		},
		Publications: []Publication{
			{
				Kind:              "publication",
				Name:              "Go Programming",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Publisher:         "Tech Publisher",
				ReleaseDate:       "2014-10-01",
				URL:               "https://publication.com",
				Summary:           "A comprehensive guide to Go programming",
			},
		},
		Skills: []Skill{
			{
				Kind:              "skill",
				Name:              "Web Development",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Level:             "Master",
				Keywords:          []string{"HTML", "CSS", "JavaScript"},
			},
		},
		Languages: []Language{
			{
				Kind:              "language",
				Name:              "English",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Fluency:           "Native speaker",
			},
		},
		Interests: []Interest{
			{
				Kind:              "interest",
				Name:              "Wildlife",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Keywords:          []string{"Ferrets", "Unicorns"},
			},
		},
		References: []Reference{
			{
				Kind:              "reference",
				Name:              "Jane Doe",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				Reference:         "Reference text",
			},
		},
		Projects: []Project{
			{
				Kind:              "project",
				Name:              "AI Research",
				Namespace:         "default",
				OwnerRef:          ownerRef,
				CreationTimestamp: timestamp,
				StartDate:         "2019-01-01",
				EndDate:           "2021-01-01",
				Description:       "Research on AI technologies",
				Highlights:        []string{"Won award at AIHacks 2016"},
				URL:               "https://project.com",
			},
		},
	}

	if resume.Kind != "resume" {
		t.Errorf("expected Kind to be 'resume', got %s", resume.Kind)
	}
	if resume.Name != "JohnDoeResume" {
		t.Errorf("expected Name to be 'JohnDoeResume', got %s", resume.Name)
	}
}
