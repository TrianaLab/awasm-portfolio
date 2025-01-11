package service

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/util"
	"fmt"

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

	kind = util.NormalizeResourceName(kind)
	logger.Trace(logrus.Fields{
		"normalized_kind": kind,
	}, "Normalized kind for GetService")

	if !util.IsValidResource(kind) {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind")
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	if name == "" {
		resources, err := s.repo.List(kind)
		if err != nil {
			logger.Error(logrus.Fields{
				"kind":  kind,
				"error": err,
			}, "Failed to list resources")
			return "", err
		}

		var filtered []models.Resource
		for _, res := range resources {
			if namespace == "" || res.GetNamespace() == namespace {
				filtered = append(filtered, res)
			}
		}

		if len(filtered) == 0 {
			logger.Info(logrus.Fields{
				"kind":      kind,
				"namespace": namespace,
			}, "No resources found")
			return "No resources found.", nil
		}

		var result string
		for _, res := range filtered {
			result += fmt.Sprintf("- %s/%s\n", kind, res.GetName())
		}

		logger.Info(logrus.Fields{
			"kind":      kind,
			"namespace": namespace,
			"count":     len(filtered),
		}, "Resources retrieved successfully")
		return result, nil
	}

	resource, err := s.repo.Get(kind, name)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
			"error":     err,
		}, "Resource not found")
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	if namespace != "" && resource.GetNamespace() != namespace {
		logger.Error(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "Resource found in a different namespace")
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	logger.Info(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "Resource retrieved successfully")
	return fmt.Sprintf("Name: %s\nNamespace: %s\n", resource.GetName(), resource.GetNamespace()), nil
}
