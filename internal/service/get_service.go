package service

import (
	"awasm-portfolio/internal/logger"
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

	if len(resources) == 0 {
		logger.Info(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "No resources found")
		return "No resources found.", nil
	}

	var result string
	for _, res := range resources {
		result += fmt.Sprintf("- %s/%s in namespace '%s'\n", res.GetKind(), res.GetName(), res.GetNamespace())
	}

	logger.Info(logrus.Fields{
		"kind":      kind,
		"namespace": namespace,
		"count":     len(resources),
	}, "Resources retrieved successfully")
	return result, nil
}
