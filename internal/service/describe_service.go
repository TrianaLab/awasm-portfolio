package service

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
	"strings"
)

type DescribeService struct {
	repo *repository.InMemoryRepository
}

func NewDescribeService(repo *repository.InMemoryRepository) *DescribeService {
	return &DescribeService{repo: repo}
}

func (s *DescribeService) DescribeResource(kind, name, namespace string) (string, error) {
	if name != "" && namespace == "" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

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

	formatter := ui.NewDetailsFormatter()

	var detailsBuilder strings.Builder
	for _, resource := range resources {
		detailsBuilder.WriteString(formatter.FormatDetails(resource))
		detailsBuilder.WriteString("---\n")
	}

	return detailsBuilder.String(), nil
}
