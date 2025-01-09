package factory

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
)

type ResourceFactory struct{}

func (f *ResourceFactory) Create(kind string, data map[string]interface{}) models.Resource {
	switch kind {
	case "profile":
		return &types.Profile{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	case "namespace":
		return &types.Namespace{Name: data["name"].(string)}
	case "education":
		return &types.Education{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	case "experience":
		return &types.Experience{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	case "contact":
		return &types.Contact{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	case "certifications":
		return &types.Certifications{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	case "contributions":
		return &types.Contributions{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	case "skills":
		return &types.Skills{Name: data["name"].(string), Namespace: data["namespace"].(string)}
	default:
		panic("Unknown resource kind")
	}
}
