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
		faker: gofakeit.New(0),
	}
}

func (f *ResourceFactory) Create(kind string, data map[string]interface{}) models.Resource {
	timestamp := time.Now()
	ownerRef := models.OwnerReference{
		Kind: "Namespace",
		Name: data["namespace"].(string),
	}

	switch kind {
	case "resume":
		return &types.Resume{
			Kind:              kind,
			Name:              data["name"].(string),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Basics: types.Basics{
				Name:    f.faker.Name(),
				Label:   f.faker.JobTitle(),
				Image:   f.faker.ImageURL(300, 300),
				Email:   f.faker.Email(),
				Phone:   f.faker.Phone(),
				Url:     f.faker.URL(),
				Summary: f.faker.Paragraph(3, 5, 10, "."),
				Location: types.Location{
					Address:     f.faker.Address().Address,
					PostalCode:  f.faker.Address().Zip,
					City:        f.faker.Address().City,
					CountryCode: f.faker.Address().Country,
					Region:      f.faker.Address().State,
				},
				Profiles: []types.Profile{
					{
						Network:  "LinkedIn",
						Username: f.faker.Username(),
						Url:      "https://linkedin.com/in/" + f.faker.Username(),
					},
					{
						Network:  "GitHub",
						Username: f.faker.Username(),
						Url:      "https://github.com/" + f.faker.Username(),
					},
				},
			},
			Work: []types.Work{
				{
					Kind:              "work",
					Name:              f.faker.Company(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Position:          f.faker.JobTitle(),
					URL:               f.faker.URL(),
					StartDate:         f.faker.Date().Format("2006-01-02"),
					EndDate:           f.faker.Date().Format("2006-01-02"),
					Summary:           f.faker.Paragraph(2, 4, 8, "."),
					Highlights:        []string{f.faker.Sentence(10), f.faker.Sentence(8)},
				},
			},
			Volunteer: []types.Volunteer{
				{
					Kind:              "volunteer",
					Name:              f.faker.Company(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Position:          "Volunteer " + f.faker.JobTitle(),
					URL:               f.faker.URL(),
					StartDate:         f.faker.Date().Format("2006-01-02"),
					EndDate:           f.faker.Date().Format("2006-01-02"),
					Summary:           f.faker.Paragraph(2, 4, 8, "."),
					Highlights:        []string{f.faker.Sentence(10)},
				},
			},
			Education: []types.Education{
				{
					Kind:              "education",
					Name:              f.faker.Company() + " University",
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					URL:               f.faker.URL(),
					Area:              f.faker.JobLevel(),
					StudyType:         "Bachelor",
					StartDate:         f.faker.Date().Format("2006-01-02"),
					EndDate:           f.faker.Date().Format("2006-01-02"),
					Score:             fmt.Sprintf("%.2f", f.faker.Float32Range(3.0, 4.0)),
					Courses:           []string{f.faker.JobDescriptor(), f.faker.JobDescriptor()},
				},
			},
			Awards: []types.Award{
				{
					Kind:              "award",
					Name:              "Best " + f.faker.JobTitle(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Date:              f.faker.Date().Format("2006-01-02"),
					Awarder:           f.faker.Company(),
					Summary:           f.faker.Sentence(15),
				},
			},
			Certificates: []types.Certificate{
				{
					Kind:              "certificate",
					Name:              f.faker.JobTitle() + " Certificate",
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Date:              f.faker.Date().Format("2006-01-02"),
					Issuer:            f.faker.Company(),
					URL:               f.faker.URL(),
				},
			},
			Publications: []types.Publication{
				{
					Kind:              "publication",
					Name:              f.faker.JobTitle() + " Research",
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Publisher:         f.faker.Company(),
					ReleaseDate:       f.faker.Date().Format("2006-01-02"),
					URL:               f.faker.URL(),
					Summary:           f.faker.Paragraph(2, 4, 8, "."),
				},
			},
			Skills: []types.Skill{
				{
					Kind:              "skill",
					Name:              f.faker.ProgrammingLanguage(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Level:             "Expert",
					Keywords:          []string{f.faker.Word(), f.faker.Word(), f.faker.Word()},
				},
			},
			Languages: []types.Language{
				{
					Kind:              "language",
					Name:              f.faker.Language(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Fluency:           "Native",
				},
			},
			Interests: []types.Interest{
				{
					Kind:              "interest",
					Name:              f.faker.Hobby(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Keywords:          []string{f.faker.Word(), f.faker.Word()},
				},
			},
			References: []types.Reference{
				{
					Kind:              "reference",
					Name:              f.faker.Name(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					Reference:         f.faker.Paragraph(1, 3, 5, "."),
				},
			},
			Projects: []types.Project{
				{
					Kind:              "project",
					Name:              f.faker.AppName(),
					Namespace:         data["namespace"].(string),
					OwnerRef:          ownerRef,
					CreationTimestamp: timestamp,
					StartDate:         f.faker.Date().Format("2006-01-02"),
					EndDate:           f.faker.Date().Format("2006-01-02"),
					Description:       f.faker.Paragraph(2, 4, 8, "."),
					Highlights:        []string{f.faker.Sentence(10)},
					URL:               f.faker.URL(),
				},
			},
		}
	case "work":
		return &types.Work{
			Kind:              kind,
			Name:              f.faker.Company(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Position:          f.faker.JobTitle(),
			URL:               f.faker.URL(),
			StartDate:         f.faker.Date().Format("2006-01-02"),
			EndDate:           f.faker.Date().Format("2006-01-02"),
			Summary:           f.faker.Paragraph(2, 4, 8, "."),
			Highlights:        []string{f.faker.Sentence(10), f.faker.Sentence(8)},
		}

	case "volunteer":
		return &types.Volunteer{
			Kind:              kind,
			Name:              f.faker.Company(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Position:          "Volunteer " + f.faker.JobTitle(),
			URL:               f.faker.URL(),
			StartDate:         f.faker.Date().Format("2006-01-02"),
			EndDate:           f.faker.Date().Format("2006-01-02"),
			Summary:           f.faker.Paragraph(2, 4, 8, "."),
			Highlights:        []string{f.faker.Sentence(10)},
		}

	case "education":
		return &types.Education{
			Kind:              kind,
			Name:              f.faker.Company() + " University",
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			URL:               f.faker.URL(),
			Area:              f.faker.JobLevel(),
			StudyType:         "Bachelor",
			StartDate:         f.faker.Date().Format("2006-01-02"),
			EndDate:           f.faker.Date().Format("2006-01-02"),
			Score:             "4.0",
			Courses:           []string{f.faker.JobDescriptor(), f.faker.JobDescriptor()},
		}

	case "award":
		return &types.Award{
			Kind:              kind,
			Name:              "Best " + f.faker.JobTitle(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Date:              f.faker.Date().Format("2006-01-02"),
			Awarder:           f.faker.Company(),
			Summary:           f.faker.Sentence(15),
		}

	case "certificate":
		return &types.Certificate{
			Kind:              kind,
			Name:              f.faker.JobTitle() + " Certificate",
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Date:              f.faker.Date().Format("2006-01-02"),
			Issuer:            f.faker.Company(),
			URL:               f.faker.URL(),
		}

	case "publication":
		return &types.Publication{
			Kind:              kind,
			Name:              f.faker.JobTitle() + " Research",
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Publisher:         f.faker.Company(),
			ReleaseDate:       f.faker.Date().Format("2006-01-02"),
			URL:               f.faker.URL(),
			Summary:           f.faker.Paragraph(2, 4, 8, "."),
		}

	case "skill":
		return &types.Skill{
			Kind:              kind,
			Name:              f.faker.ProgrammingLanguage(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Level:             "Expert",
			Keywords:          []string{f.faker.Word(), f.faker.Word(), f.faker.Word()},
		}

	case "language":
		return &types.Language{
			Kind:              kind,
			Name:              f.faker.Language(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Fluency:           "Native speaker",
		}

	case "interest":
		return &types.Interest{
			Kind:              kind,
			Name:              f.faker.Hobby(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Keywords:          []string{f.faker.Word(), f.faker.Word()},
		}

	case "reference":
		return &types.Reference{
			Kind:              kind,
			Name:              f.faker.Name(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			Reference:         f.faker.Paragraph(1, 3, 5, "."),
		}

	case "project":
		return &types.Project{
			Kind:              kind,
			Name:              f.faker.AppName(),
			Namespace:         data["namespace"].(string),
			OwnerRef:          ownerRef,
			CreationTimestamp: timestamp,
			StartDate:         f.faker.Date().Format("2006-01-02"),
			EndDate:           f.faker.Date().Format("2006-01-02"),
			Description:       f.faker.Paragraph(2, 4, 8, "."),
			Highlights:        []string{f.faker.Sentence(10)},
			URL:               f.faker.URL(),
		}

	default:
		return nil
	}
}
