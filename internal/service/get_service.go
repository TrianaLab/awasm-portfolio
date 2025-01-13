package service

import (
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
	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		return "", err
	}

	// Apply namespace-specific logic when "all" is requested
	if strings.ToLower(kind) == "all" && namespace != "" {
		for i := 0; i < len(resources); {
			if resources[i].GetKind() == "namespace" && resources[i].GetName() != namespace {
				resources = append(resources[:i], resources[i+1:]...)
			} else {
				i++
			}
		}
	}

	// Instantiate the UnifiedFormatter and format the table
	formatter := ui.NewUnifiedFormatter()
	return formatter.FormatTable(resources), nil
}
