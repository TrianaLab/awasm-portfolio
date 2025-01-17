package service

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
	"strings"
)

type GetService struct {
	repo *repository.InMemoryRepository
}

func NewGetService(repo *repository.InMemoryRepository) *GetService {
	return &GetService{repo: repo}
}

func (s *GetService) GetResources(kind string, name string, namespace string) (string, error) {
	if name != "" && namespace == "" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

	// Retrieve resources from the repository
	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		return "", err
	}

	// Filter namespaces when kind is "all" and namespace is specified
	if strings.ToLower(kind) == "all" && namespace != "" {
		for i := 0; i < len(resources); {
			if resources[i].GetKind() == "namespace" && resources[i].GetName() != namespace {
				resources = append(resources[:i], resources[i+1:]...)
			} else {
				i++
			}
		}
	}

	// Instantiate and use the new TableFormatter
	formatter := ui.NewTableFormatter()
	return formatter.FormatTable(resources), nil
}
