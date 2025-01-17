package ui

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"fmt"
	"reflect"
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
		"profile": {
			Headers: []string{"NAME", "NAMESPACE", "CONTRIBUTIONS", "EXPERIENCE", "CERTIFICATIONS", "EDUCATION", "SKILLS", "CONTACT", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					if profile, ok := r.(*types.Profile); ok {
						return profile.Contributions.GetName()
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if profile, ok := r.(*types.Profile); ok {
						return profile.Experience.GetName()
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if profile, ok := r.(*types.Profile); ok {
						return profile.Certifications.GetName()
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if profile, ok := r.(*types.Profile); ok {
						return profile.Education.GetName()
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if profile, ok := r.(*types.Profile); ok {
						return profile.Skills.GetName()
					}
					return "N/A"
				},
				func(r models.Resource) string {
					if profile, ok := r.(*types.Profile); ok {
						if profile.Contact.Email != "" {
							return profile.Contact.GetName()
						}
					}
					return "N/A"
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
		"default": {
			Headers: []string{"NAME", "NAMESPACE", "ITEMS", "AGE"},
			Extractors: []func(models.Resource) string{
				func(r models.Resource) string { return r.GetName() },
				func(r models.Resource) string { return r.GetNamespace() },
				func(r models.Resource) string {
					// Use reflection to count slice fields dynamically
					val := reflect.ValueOf(r).Elem()
					count := 0
					for i := 0; i < val.NumField(); i++ {
						field := val.Field(i)
						if field.Kind() == reflect.Slice {
							count += field.Len()
						}
					}
					return fmt.Sprintf("%d items", count)
				},
				func(r models.Resource) string { return calculateAge(r.GetCreationTimestamp()) },
			},
		},
	}
}
