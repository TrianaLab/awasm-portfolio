package ui

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"fmt"
)

// Schema defines headers and extractors for resources
type Schema struct {
	Headers    []string
	Extractors []func(models.Resource) string
}

// GenerateSchemas creates schemas for all resource types dynamically
func GenerateSchemas() map[string]Schema {
	return map[string]Schema{
		"namespace": {
			Headers: []string{"NAME", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"resume": {
			Headers: []string{"NAME", "NAMESPACE", "BASICS", "WORK", "VOLUNTEER", "EDUCATION", "AWARDS", "CERTIFICATES", "PUBLICATIONS", "SKILLS", "LANGUAGES", "INTERESTS", "REFERENCES", "PROJECTS", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return resume.Basics.Name
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Work))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Volunteer))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Education))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Awards))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Certificates))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Publications))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Skills))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Languages))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Interests))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.References))
					}
					return "0"
				},
				func(r models.Resource) string {
					if resume, ok := r.(*types.Resume); ok {
						return fmt.Sprintf("%d", len(resume.Projects))
					}
					return "0"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"basics": {
			Headers: []string{"NAME", "NAMESPACE", "FULL NAME", "LABEL", "EMAIL", "PHONE", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if basics, ok := r.(*types.Basics); ok {
						return basics.FullName
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if basics, ok := r.(*types.Basics); ok {
						return basics.Label
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if basics, ok := r.(*types.Basics); ok {
						return basics.Email
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if basics, ok := r.(*types.Basics); ok {
						return basics.Phone
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"work": {
			Headers: []string{"NAME", "NAMESPACE", "COMPANY", "POSITION", "START", "END", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if work, ok := r.(*types.Work); ok {
						return work.Company
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if work, ok := r.(*types.Work); ok {
						return work.Position
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if work, ok := r.(*types.Work); ok {
						return work.StartDate
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if work, ok := r.(*types.Work); ok {
						if work.EndDate == "" {
							return "Present"
						}
						return work.EndDate
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"volunteer": {
			Headers: []string{"NAME", "NAMESPACE", "ORGANIZATION", "POSITION", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if vol, ok := r.(*types.Volunteer); ok {
						return vol.Organization
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if vol, ok := r.(*types.Volunteer); ok {
						return vol.Position
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"education": {
			Headers: []string{"NAME", "NAMESPACE", "INSTITUTION", "AREA", "STUDY TYPE", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if edu, ok := r.(*types.Education); ok {
						return edu.Institution
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if edu, ok := r.(*types.Education); ok {
						return edu.Area
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if edu, ok := r.(*types.Education); ok {
						return edu.StudyType
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"skill": {
			Headers: []string{"NAME", "NAMESPACE", "SKILL", "LEVEL", "KEYWORDS", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if skill, ok := r.(*types.Skill); ok {
						return skill.Skill
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if skill, ok := r.(*types.Skill); ok {
						return skill.Level
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if skill, ok := r.(*types.Skill); ok {
						if len(skill.Keywords) > 3 {
							return fmt.Sprintf("%s, ...", skill.Keywords[0:3])
						}
						return fmt.Sprintf("%v", skill.Keywords)
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"language": {
			Headers: []string{"NAME", "NAMESPACE", "LANGUAGE", "FLUENCY", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if lang, ok := r.(*types.Language); ok {
						return lang.Language
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if lang, ok := r.(*types.Language); ok {
						return lang.Fluency
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"project": {
			Headers: []string{"NAME", "NAMESPACE", "PROJECT", "URL", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if proj, ok := r.(*types.Project); ok {
						return proj.Project
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if proj, ok := r.(*types.Project); ok {
						return proj.URL
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"publication": {
			Headers: []string{"NAME", "NAMESPACE", "PUBLICATION", "PUBLISHER", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if pub, ok := r.(*types.Publication); ok {
						return pub.Publication
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if pub, ok := r.(*types.Publication); ok {
						return pub.Publisher
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"certificate": {
			Headers: []string{"NAME", "NAMESPACE", "CERTIFICATE", "DATE", "ISSUER", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if certificate, ok := r.(*types.Certificate); ok {
						return certificate.Certificate
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if certificate, ok := r.(*types.Certificate); ok {
						return certificate.Date
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if certificate, ok := r.(*types.Certificate); ok {
						return certificate.Issuer
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"interest": {
			Headers: []string{"NAME", "NAMESPACE", "INTEREST", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if interest, ok := r.(*types.Interest); ok {
						return interest.Interest
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"award": {
			Headers: []string{"NAME", "NAMESPACE", "TITLE", "AWARDER", "DATE", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if award, ok := r.(*types.Award); ok {
						return award.Title
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if award, ok := r.(*types.Award); ok {
						return award.Awarder
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if award, ok := r.(*types.Award); ok {
						return award.Date
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"default": {
			Headers: []string{"NAME", "NAMESPACE", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
	}
}
