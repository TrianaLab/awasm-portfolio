package service

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"awasm-portfolio/internal/util"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type DescribeService struct {
	repo *repository.InMemoryRepository
}

func NewDescribeService(repo *repository.InMemoryRepository) *DescribeService {
	return &DescribeService{repo: repo}
}

func (s *DescribeService) DescribeResource(kind, name, namespace string) (string, error) {
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind in List")
		return "", err
	}

	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve resources: %w", err)
	}

	if len(resources) == 0 {
		return "No resources found.", nil
	}

	// Instantiate the UnifiedFormatter for formatting details
	formatter := ui.NewUnifiedFormatter()

	var detailsBuilder strings.Builder
	for _, resource := range resources {
		detailsBuilder.WriteString(formatter.FormatDetails(resource))
		detailsBuilder.WriteString("\n---\n")
	}

	return detailsBuilder.String(), nil
}
