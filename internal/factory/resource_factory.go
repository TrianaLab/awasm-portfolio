package factory

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type ResourceFactory struct {
	faker *gofakeit.Faker
}

func NewResourceFactory() *ResourceFactory {
	return &ResourceFactory{
		faker: gofakeit.New(time.Now().Unix()),
	}
}

func (f *ResourceFactory) Create(kind string, data map[string]interface{}) models.Resource {
	switch kind {
	case "profile":
		return &types.Profile{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Certifications: types.Certifications{
				Kind:      kind,
				Name:      data["name"].(string),
				Namespace: data["namespace"].(string),
				Certifications: []types.Certification{
					{Description: f.faker.JobTitle() + " Certification", Link: f.faker.URL()},
					{Description: f.faker.JobTitle() + " Advanced Certification", Link: f.faker.URL()},
				},
			},
			Contact: types.Contact{
				Kind:      kind,
				Name:      data["name"].(string),
				Namespace: data["namespace"].(string),
				Email:     f.faker.Email(),
				Linkedin:  "https://linkedin.com/in/" + f.faker.Username(),
				Github:    "https://github.com/" + f.faker.Username(),
			},
			Contributions: types.Contributions{
				Kind:      kind,
				Name:      data["name"].(string),
				Namespace: data["namespace"].(string),
				Contributions: []types.Contribution{
					{Project: f.faker.BuzzWord(), Description: f.faker.HackerPhrase(), Link: f.faker.URL()},
					{Project: f.faker.BuzzWord(), Description: f.faker.HackerPhrase(), Link: f.faker.URL()},
				},
			},
			Education: types.Education{
				Kind:      kind,
				Name:      data["name"].(string),
				Namespace: data["namespace"].(string),
				Courses: []types.Course{
					{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
					{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
					{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
					{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
				},
			},
			Experience: types.Experience{
				Kind:      kind,
				Name:      data["name"].(string),
				Namespace: data["namespace"].(string),
				Jobs: []types.Job{
					{
						Title:       f.faker.JobTitle(),
						Description: f.faker.JobDescriptor(),
						Company:     f.faker.Company(),
						Duration:    fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year()),
					},
					{
						Title:       f.faker.JobTitle(),
						Description: f.faker.JobDescriptor(),
						Company:     f.faker.Company(),
						Duration:    fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year()),
					},
				},
			},
			Skills: types.Skills{
				Kind:      kind,
				Name:      data["name"].(string),
				Namespace: data["namespace"].(string),
				Skills: []types.Skill{
					{Category: "Programming", Competence: f.faker.ProgrammingLanguage(), Proficiency: "Expert"},
					{Category: "Programming", Competence: f.faker.ProgrammingLanguage(), Proficiency: "Expert"},
					{Category: "DevOps", Competence: f.faker.BuzzWord(), Proficiency: "Advanced"},
					{Category: "Languages", Competence: f.faker.Language(), Proficiency: "Fluent"},
					{Category: "Languages", Competence: f.faker.Language(), Proficiency: "Native"},
				},
			},
		}
	case "namespace":
		return &types.Namespace{
			Name: data["name"].(string),
			Kind: "namespace",
		}
	case "education":
		return &types.Education{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Courses: []types.Course{
				{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
				{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
				{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
				{Title: f.faker.JobTitle(), Institution: f.faker.Company(), Duration: fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year())},
			},
		}
	case "experience":
		return &types.Experience{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Jobs: []types.Job{
				{
					Title:       f.faker.JobTitle(),
					Description: f.faker.JobDescriptor(),
					Company:     f.faker.Company(),
					Duration:    fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year()),
				},
				{
					Title:       f.faker.JobTitle(),
					Description: f.faker.JobDescriptor(),
					Company:     f.faker.Company(),
					Duration:    fmt.Sprintf("%d - %d", f.faker.Year(), f.faker.Year()),
				},
			},
		}
	case "contact":
		return &types.Contact{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Email:     f.faker.Email(),
			Linkedin:  "https://linkedin.com/in/" + f.faker.Username(),
			Github:    "https://github.com/" + f.faker.Username(),
		}
	case "certifications":
		return &types.Certifications{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Certifications: []types.Certification{
				{Description: f.faker.JobTitle() + " Certification", Link: f.faker.URL()},
				{Description: f.faker.JobTitle() + " Advanced Certification", Link: f.faker.URL()},
			},
		}
	case "contributions":
		return &types.Contributions{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Contributions: []types.Contribution{
				{Project: f.faker.BuzzWord(), Description: f.faker.HackerPhrase(), Link: f.faker.URL()},
				{Project: f.faker.BuzzWord(), Description: f.faker.HackerPhrase(), Link: f.faker.URL()},
			},
		}
	case "skills":
		return &types.Skills{
			Kind:      kind,
			Name:      data["name"].(string),
			Namespace: data["namespace"].(string),
			Skills: []types.Skill{
				{Category: "Programming", Competence: f.faker.ProgrammingLanguage(), Proficiency: "Expert"},
				{Category: "Programming", Competence: f.faker.ProgrammingLanguage(), Proficiency: "Expert"},
				{Category: "DevOps", Competence: f.faker.BuzzWord(), Proficiency: "Advanced"},
				{Category: "Languages", Competence: f.faker.Language(), Proficiency: "Fluent"},
				{Category: "Languages", Competence: f.faker.Language(), Proficiency: "Native"},
			},
		}
	default:
		return nil
	}
}
