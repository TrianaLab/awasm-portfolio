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
