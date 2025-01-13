package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"strings"
)

type GetService struct {
	repo *repository.InMemoryRepository
}

func NewGetService(repo *repository.InMemoryRepository) *GetService {
	return &GetService{repo: repo}
}

func (s *GetService) GetResources(kind string, name string, namespace string) (string, error) {
	all := strings.ToLower(kind) == "all"
	// Retrieve resources from the repository
	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		return "", err
	}

	// Apply namespace-specific logic when "all" is requested
	if all && namespace != "" {
		filteredResources := []models.Resource{}
		for _, res := range resources {
			if res.GetKind() == "namespace" && res.GetName() == namespace {
				filteredResources = append(filteredResources, res)
			} else if res.GetNamespace() == namespace {
				filteredResources = append(filteredResources, res)
			}
		}
		resources = filteredResources
	}

	if len(resources) == 0 {
		return "No resources found.", nil
	}

	// Instantiate the UnifiedFormatter and format the table
	formatter := ui.NewUnifiedFormatter()
	return formatter.FormatTable(resources), nil
}
