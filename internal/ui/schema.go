package ui

import (
	"fmt"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
)

// Schema defines headers and extractors for resources
type Schema struct {
	Headers    []string
	Extractors []func(models.Resource) string
}

// GenerateSchemas creates schemas for all resource types dynamically.
func GenerateSchemas() map[string]Schema {
	return map[string]Schema{
		"namespace":   namespaceSchema(),
		"resume":      resumeSchema(),
		"basics":      basicsSchema(),
		"work":        workSchema(),
		"volunteer":   volunteerSchema(),
		"education":   educationSchema(),
		"skill":       skillSchema(),
		"language":    languageSchema(),
		"project":     projectSchema(),
		"publication": publicationSchema(),
		"certificate": certificateSchema(),
		"interest":    interestSchema(),
		"award":       awardSchema(),
		"default":     defaultSchema(),
	}
}

// Per-type extractor wrappers. Each performs the type assertion and returns
// "N/A" when the resource is not of the expected concrete type. Generics
// cannot replace these because *T type assertions on an interface value are
// not legal in Go.

func resumeExtract(pick func(*types.Resume) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Resume); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func basicsExtract(pick func(*types.Basics) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Basics); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func workExtract(pick func(*types.Work) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Work); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func volunteerExtract(pick func(*types.Volunteer) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Volunteer); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func educationExtract(pick func(*types.Education) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Education); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func skillExtract(pick func(*types.Skill) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Skill); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func languageExtract(pick func(*types.Language) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Language); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func projectExtract(pick func(*types.Project) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Project); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func publicationExtract(pick func(*types.Publication) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Publication); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func certificateExtract(pick func(*types.Certificate) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Certificate); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func interestExtract(pick func(*types.Interest) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Interest); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func awardExtract(pick func(*types.Award) string) func(models.Resource) string {
	return func(r models.Resource) string {
		if v, ok := r.(*types.Award); ok {
			return pick(v)
		}
		return "N/A"
	}
}

func nameExtractor() func(models.Resource) string {
	return func(r models.Resource) string { return r.GetName() }
}

func namespaceExtractor() func(models.Resource) string {
	return func(r models.Resource) string { return r.GetNamespace() }
}

func ageExtractor() func(models.Resource) string {
	return func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) }
}

func namespaceSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			ageExtractor(),
		},
	}
}

func defaultSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			ageExtractor(),
		},
	}
}

func resumeSchema() Schema {
	count := func(pick func(*types.Resume) int) func(models.Resource) string {
		return resumeExtract(func(r *types.Resume) string { return fmt.Sprintf("%d", pick(r)) })
	}
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "BASICS", "WORK", "VOLUNTEER", "EDUCATION", "AWARDS", "CERTIFICATES", "PUBLICATIONS", "SKILLS", "LANGUAGES", "INTERESTS", "REFERENCES", "PROJECTS", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			resumeExtract(func(r *types.Resume) string { return r.Basics.Name }),
			count(func(r *types.Resume) int { return len(r.Work) }),
			count(func(r *types.Resume) int { return len(r.Volunteer) }),
			count(func(r *types.Resume) int { return len(r.Education) }),
			count(func(r *types.Resume) int { return len(r.Awards) }),
			count(func(r *types.Resume) int { return len(r.Certificates) }),
			count(func(r *types.Resume) int { return len(r.Publications) }),
			count(func(r *types.Resume) int { return len(r.Skills) }),
			count(func(r *types.Resume) int { return len(r.Languages) }),
			count(func(r *types.Resume) int { return len(r.Interests) }),
			count(func(r *types.Resume) int { return len(r.References) }),
			count(func(r *types.Resume) int { return len(r.Projects) }),
			ageExtractor(),
		},
	}
}

func basicsSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "FULL NAME", "LABEL", "EMAIL", "PHONE", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			basicsExtract(func(b *types.Basics) string { return b.FullName }),
			basicsExtract(func(b *types.Basics) string { return b.Label }),
			basicsExtract(func(b *types.Basics) string { return b.Email }),
			basicsExtract(func(b *types.Basics) string { return b.Phone }),
			ageExtractor(),
		},
	}
}

func workSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "COMPANY", "POSITION", "START", "END", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			workExtract(func(w *types.Work) string { return w.Company }),
			workExtract(func(w *types.Work) string { return w.Position }),
			workExtract(func(w *types.Work) string { return w.StartDate }),
			workExtract(func(w *types.Work) string {
				if w.EndDate == "" {
					return "Present"
				}
				return w.EndDate
			}),
			ageExtractor(),
		},
	}
}

func volunteerSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "ORGANIZATION", "POSITION", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			volunteerExtract(func(v *types.Volunteer) string { return v.Organization }),
			volunteerExtract(func(v *types.Volunteer) string { return v.Position }),
			ageExtractor(),
		},
	}
}

func educationSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "INSTITUTION", "AREA", "STUDY TYPE", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			educationExtract(func(e *types.Education) string { return e.Institution }),
			educationExtract(func(e *types.Education) string { return e.Area }),
			educationExtract(func(e *types.Education) string { return e.StudyType }),
			ageExtractor(),
		},
	}
}

func skillSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "SKILL", "LEVEL", "KEYWORDS", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			skillExtract(func(s *types.Skill) string { return s.Skill }),
			skillExtract(func(s *types.Skill) string { return s.Level }),
			skillExtract(func(s *types.Skill) string {
				if len(s.Keywords) > 3 {
					return fmt.Sprintf("%s, ...", s.Keywords[0:3])
				}
				return fmt.Sprintf("%v", s.Keywords)
			}),
			ageExtractor(),
		},
	}
}

func languageSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "LANGUAGE", "FLUENCY", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			languageExtract(func(l *types.Language) string { return l.Language }),
			languageExtract(func(l *types.Language) string { return l.Fluency }),
			ageExtractor(),
		},
	}
}

func projectSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "PROJECT", "URL", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			projectExtract(func(p *types.Project) string { return p.Project }),
			projectExtract(func(p *types.Project) string { return p.URL }),
			ageExtractor(),
		},
	}
}

func publicationSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "PUBLICATION", "PUBLISHER", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			publicationExtract(func(p *types.Publication) string { return p.Publication }),
			publicationExtract(func(p *types.Publication) string { return p.Publisher }),
			ageExtractor(),
		},
	}
}

func certificateSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "CERTIFICATE", "DATE", "ISSUER", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			certificateExtract(func(c *types.Certificate) string { return c.Certificate }),
			certificateExtract(func(c *types.Certificate) string { return c.Date }),
			certificateExtract(func(c *types.Certificate) string { return c.Issuer }),
			ageExtractor(),
		},
	}
}

func interestSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "INTEREST", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			interestExtract(func(i *types.Interest) string { return i.Interest }),
			ageExtractor(),
		},
	}
}

func awardSchema() Schema {
	return Schema{
		Headers: []string{"NAME", "NAMESPACE", "TITLE", "AWARDER", "DATE", "AGE"},
		Extractors: []func(models.Resource) string{
			nameExtractor(),
			namespaceExtractor(),
			awardExtract(func(a *types.Award) string { return a.Title }),
			awardExtract(func(a *types.Award) string { return a.Awarder }),
			awardExtract(func(a *types.Award) string { return a.Date }),
			ageExtractor(),
		},
	}
}
