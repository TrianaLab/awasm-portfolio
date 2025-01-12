package service

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"awasm-portfolio/internal/util"

	"github.com/sirupsen/logrus"
)

type GetService struct {
	repo *repository.InMemoryRepository
}

func NewGetService(repo *repository.InMemoryRepository) *GetService {
	return &GetService{repo: repo}
}

func (s *GetService) GetResources(kind string, name string, namespace string) (string, error) {
	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "GetService.GetResources called")

	// Preserve the original kind to handle "all" explicitly
	originalKind := kind
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind in List")
		return "", err
	}

	// Retrieve resources from the repository
	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
			"error":     err,
		}, "Failed to list resources")
		return "", err
	}

	// Apply namespace-specific logic when "all" is requested
	if originalKind == "all" && namespace != "" {
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
		logger.Info(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "No resources found")
		return "No resources found.", nil
	}

	// Instantiate the UnifiedFormatter and format the table
	formatter := ui.NewUnifiedFormatter()
	return formatter.FormatTable(resources), nil
}
