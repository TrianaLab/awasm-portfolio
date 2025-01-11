package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"reflect"
	"strings"
)

func SupportedResources() map[string]string {
	supported := map[string]string{}

	types := []models.Resource{
		&types.Profile{}, &types.Namespace{}, &types.Education{},
		&types.Experience{}, &types.Contact{}, &types.Certifications{},
		&types.Contributions{}, &types.Skills{},
	}

	for _, t := range types {
		kind := normalizeKind(reflect.TypeOf(t).Elem().Name())
		supported[kind] = kind
		supported[kind+"s"] = kind
	}

	supported["ns"] = "namespace"
	return supported
}

func NormalizeResourceName(resource string) string {
	supported := SupportedResources()
	if singular, exists := supported[resource]; exists {
		return singular
	}
	return resource
}

func IsValidResource(resource string) bool {
	_, exists := SupportedResources()[resource]
	return exists
}

func normalizeKind(kind string) string {
	return strings.ToLower(kind)
}
