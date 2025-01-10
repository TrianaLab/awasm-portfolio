package service

func NormalizeResourceName(resource string) string {
	// Map plural and alias forms to singular names
	plurals := map[string]string{
		"profiles":       "profile",
		"experiences":    "experience",
		"certifications": "certification",
		"contributions":  "contribution",
		"educations":     "education",
		"skills":         "skill",
		"contacts":       "contact",
		"namespaces":     "namespace",
		"ns":             "namespace", // Alias
	}

	// Check if the resource is plural or an alias
	if singular, exists := plurals[resource]; exists {
		return singular
	}

	// Return the resource unchanged if no match is found
	return resource
}

// SupportedResources returns a list of valid resource types
func SupportedResources() map[string]struct{} {
	return map[string]struct{}{
		"profile":       {},
		"namespace":     {},
		"education":     {},
		"experience":    {},
		"contact":       {},
		"certification": {},
		"contribution":  {},
		"skill":         {},
	}
}

// IsValidResource checks if a given resource type is supported
func IsValidResource(resource string) bool {
	_, exists := SupportedResources()[resource]
	return exists
}
