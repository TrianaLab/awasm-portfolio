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
func (s *GetService) GetResources(kind, name, namespace string) (string, error) {
	if kind == "" {
		return "", fmt.Errorf("resource kind is required")
	}

	// If name is empty, list all resources of the kind
	if name == "" {
		resources := s.repo.List(kind)

		// Filter by namespace if specified
		var filtered []models.Resource
		for _, res := range resources {
			if namespace == "" || res.GetNamespace() == namespace {
				filtered = append(filtered, res)
			}
		}

		// Format the results
		if len(filtered) == 0 {
			return "No resources found.", nil
		}

		var result string
		for _, res := range filtered {
			result += fmt.Sprintf("- %s/%s\n", kind, res.GetName())
		}

		return result, nil
	}

	// If name is provided, retrieve a single resource
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	// Ensure the namespace matches
	if namespace != "" && resource.GetNamespace() != namespace {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	// Format the result for a single resource
	return fmt.Sprintf("Name: %s\nNamespace: %s\n", resource.GetName(), resource.GetNamespace()), nil
}
