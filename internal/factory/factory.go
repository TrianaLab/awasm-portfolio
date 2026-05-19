// Package factory builds populated resource instances for `kubectl
// create`. Demo data (gofakeit) is used so freshly created resumes
// look like a real portfolio. The package is small and explicit —
// one builder per kind, no reflection.
package factory

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
)

// New returns a resource of the given kind populated with plausible
// fake data. For Resume it includes child Basics + a few entries in
// every slice section so the resume view and PDF show non-empty
// content immediately.
func New(kind, name, namespace string) (models.Resource, error) {
	meta := models.Meta{Kind: kind, Name: name, Namespace: namespace}
	switch kind {
	case "namespace":
		return &types.Namespace{Meta: models.Meta{Kind: kind, Name: name}}, nil
	case "resume":
		return newResume(meta, name, namespace), nil
	case "basics":
		return newBasics(meta), nil
	case "work":
		return newWork(meta), nil
	case "volunteer":
		return newVolunteer(meta), nil
	case "education":
		return newEducation(meta), nil
	case "award":
		return newAward(meta), nil
	case "certificate":
		return newCertificate(meta), nil
	case "publication":
		return newPublication(meta), nil
	case "skill":
		return newSkill(meta), nil
	case "language":
		return newLanguage(meta), nil
	case "interest":
		return newInterest(meta), nil
	case "reference":
		return newReference(meta), nil
	case "project":
		return newProject(meta), nil
	default:
		return nil, fmt.Errorf("unsupported resource kind: %s", kind)
	}
}

func newBasics(meta models.Meta) *types.Basics {
	return &types.Basics{
		Meta:     meta,
		FullName: gofakeit.Name(),
		Label:    gofakeit.JobTitle(),
		Email:    gofakeit.Email(),
		Phone:    gofakeit.Phone(),
		Url:      gofakeit.URL(),
		Summary:  gofakeit.Sentence(20),
		Location: types.Location{
			City:        gofakeit.City(),
			Region:      gofakeit.State(),
			CountryCode: gofakeit.CountryAbr(),
			PostalCode:  gofakeit.Zip(),
		},
		Profiles: []types.Profile{
			{Network: "GitHub", Username: gofakeit.Username(), Url: gofakeit.URL()},
			{Network: "LinkedIn", Username: gofakeit.Username(), Url: gofakeit.URL()},
		},
	}
}

func newWork(meta models.Meta) *types.Work {
	return &types.Work{
		Meta:      meta,
		Company:   gofakeit.Company(),
		Position:  gofakeit.JobTitle(),
		URL:       gofakeit.URL(),
		StartDate: gofakeit.Date().Format("2006-01-02"),
		EndDate:   gofakeit.Date().Format("2006-01-02"),
		Summary:   gofakeit.Sentence(15),
	}
}

func newVolunteer(meta models.Meta) *types.Volunteer {
	return &types.Volunteer{
		Meta:         meta,
		Organization: gofakeit.Company(),
		Position:     gofakeit.JobTitle(),
		URL:          gofakeit.URL(),
		StartDate:    gofakeit.Date().Format("2006-01-02"),
		EndDate:      gofakeit.Date().Format("2006-01-02"),
		Summary:      gofakeit.Sentence(15),
	}
}

func newEducation(meta models.Meta) *types.Education {
	return &types.Education{
		Meta:        meta,
		Institution: gofakeit.Company() + " University",
		URL:         gofakeit.URL(),
		Area:        gofakeit.JobDescriptor(),
		StudyType:   "B.Eng",
		StartDate:   gofakeit.Date().Format("2006-01-02"),
		EndDate:     gofakeit.Date().Format("2006-01-02"),
		Courses:     []string{gofakeit.BuzzWord(), gofakeit.BuzzWord(), gofakeit.BuzzWord()},
	}
}

func newAward(meta models.Meta) *types.Award {
	return &types.Award{
		Meta:    meta,
		Title:   gofakeit.BuzzWord() + " Award",
		Date:    gofakeit.Date().Format("2006-01-02"),
		Awarder: gofakeit.Company(),
		Summary: gofakeit.Sentence(10),
	}
}

func newCertificate(meta models.Meta) *types.Certificate {
	return &types.Certificate{
		Meta:        meta,
		Certificate: gofakeit.BuzzWord() + " Certified",
		Date:        gofakeit.Date().Format("2006-01-02"),
		Issuer:      gofakeit.Company(),
		URL:         gofakeit.URL(),
	}
}

func newPublication(meta models.Meta) *types.Publication {
	return &types.Publication{
		Meta:        meta,
		Publication: gofakeit.HipsterSentence(5),
		Publisher:   gofakeit.Company(),
		ReleaseDate: gofakeit.Date().Format("2006-01-02"),
		URL:         gofakeit.URL(),
		Summary:     gofakeit.Sentence(15),
	}
}

func newSkill(meta models.Meta) *types.Skill {
	return &types.Skill{
		Meta:     meta,
		Skill:    gofakeit.BuzzWord(),
		Level:    gofakeit.RandomString([]string{"Beginner", "Intermediate", "Advanced", "Expert"}),
		Keywords: []string{gofakeit.HackerNoun(), gofakeit.HackerNoun(), gofakeit.HackerNoun()},
	}
}

func newLanguage(meta models.Meta) *types.Language {
	return &types.Language{
		Meta:     meta,
		Language: gofakeit.Language(),
		Fluency:  gofakeit.RandomString([]string{"Basic", "Conversational", "Fluent", "Native"}),
	}
}

func newInterest(meta models.Meta) *types.Interest {
	return &types.Interest{
		Meta:     meta,
		Interest: gofakeit.Hobby(),
		Keywords: []string{gofakeit.BuzzWord(), gofakeit.BuzzWord()},
	}
}

func newReference(meta models.Meta) *types.Reference {
	return &types.Reference{
		Meta:      meta,
		Person:    gofakeit.Name(),
		Reference: gofakeit.Sentence(15),
	}
}

func newProject(meta models.Meta) *types.Project {
	return &types.Project{
		Meta:        meta,
		Project:     gofakeit.AppName(),
		StartDate:   gofakeit.Date().Format("2006-01-02"),
		EndDate:     gofakeit.Date().Format("2006-01-02"),
		Description: gofakeit.Sentence(15),
		URL:         gofakeit.URL(),
		Highlights:  []string{gofakeit.BuzzWord(), gofakeit.BuzzWord()},
	}
}

// newResume builds a Resume populated with a Basics block + 2 entries
// in every slice section. The child resources also get sensible names
// derived from the resume name so kubectl get <kind> shows them.
func newResume(meta models.Meta, name, namespace string) *types.Resume {
	r := &types.Resume{
		Meta:   meta,
		Basics: *newBasics(childMeta("basics", name+"-basics", namespace, name)),
	}
	for i := range 2 {
		r.Work = append(r.Work, *newWork(childMeta("work", fmt.Sprintf("%s-work-%d", name, i), namespace, name)))
		r.Volunteer = append(r.Volunteer, *newVolunteer(childMeta("volunteer", fmt.Sprintf("%s-volunteer-%d", name, i), namespace, name)))
		r.Education = append(r.Education, *newEducation(childMeta("education", fmt.Sprintf("%s-education-%d", name, i), namespace, name)))
		r.Awards = append(r.Awards, *newAward(childMeta("award", fmt.Sprintf("%s-award-%d", name, i), namespace, name)))
		r.Certificates = append(r.Certificates, *newCertificate(childMeta("certificate", fmt.Sprintf("%s-certificate-%d", name, i), namespace, name)))
		r.Publications = append(r.Publications, *newPublication(childMeta("publication", fmt.Sprintf("%s-publication-%d", name, i), namespace, name)))
		r.Skills = append(r.Skills, *newSkill(childMeta("skill", fmt.Sprintf("%s-skill-%d", name, i), namespace, name)))
		r.Languages = append(r.Languages, *newLanguage(childMeta("language", fmt.Sprintf("%s-language-%d", name, i), namespace, name)))
		r.Interests = append(r.Interests, *newInterest(childMeta("interest", fmt.Sprintf("%s-interest-%d", name, i), namespace, name)))
		r.References = append(r.References, *newReference(childMeta("reference", fmt.Sprintf("%s-reference-%d", name, i), namespace, name)))
		r.Projects = append(r.Projects, *newProject(childMeta("project", fmt.Sprintf("%s-project-%d", name, i), namespace, name)))
	}
	return r
}

func childMeta(kind, name, namespace, parentName string) models.Meta {
	return models.Meta{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
		OwnerRef:  models.OwnerReference{Kind: "resume", Name: parentName, Namespace: namespace},
	}
}
