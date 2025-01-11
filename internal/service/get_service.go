package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"fmt"
)

type GetService struct {
	repo *repository.InMemoryRepository
}

func NewGetService(repo *repository.InMemoryRepository) *GetService {
	return &GetService{repo: repo}
}

// GetResources retrieves resources or a single resource by kind, name, and namespace
func (s *GetService) GetResources(kind string, name string, namespace string) (string, error) {
	kind = NormalizeResourceName(kind)
	if !IsValidResource(kind) {
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	if name == "" {
		resources, err := s.repo.List(kind)
		if err != nil {
			return "", err
		}

		var filtered []models.Resource
		for _, res := range resources {
			if namespace == "" || res.GetNamespace() == namespace {
				filtered = append(filtered, res)
			}
		}

		if len(filtered) == 0 {
			return "No resources found.", nil
		}

		var result string
		for _, res := range filtered {
			result += fmt.Sprintf("- %s/%s\n", kind, res.GetName())
		}

		return result, nil
	}

	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	if namespace != "" && resource.GetNamespace() != namespace {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	return fmt.Sprintf("Name: %s\nNamespace: %s\n", resource.GetName(), resource.GetNamespace()), nil
}
