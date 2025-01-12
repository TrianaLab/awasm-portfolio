package service

import (
	"awasm-portfolio/internal/logger"
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

	kind, err := util.NormalizeKind(kind)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind in List")
		return "", err
	}

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

	// Use TextFormatter to format the resources into a table
	formatter := ui.TextFormatter{}
	output := formatter.FormatTable(resources)

	return output, nil
}
