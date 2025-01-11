package service

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"

	"github.com/sirupsen/logrus"
)

type DescribeService struct {
	repo      *repository.InMemoryRepository
	formatter ui.Formatter
}

func NewDescribeService(repo *repository.InMemoryRepository) *DescribeService {
	return &DescribeService{
		repo:      repo,
		formatter: ui.TextFormatter{},
	}
}

func (s *DescribeService) DescribeResource(kind string, name string, namespace string) (string, error) {
	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "DescribeService.DescribeResource called")

	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Failed to describe resource")
		return "", err
	}

	if len(resources) == 0 {
		logger.Info(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "No resources found for description")
		return "No resources found.", nil
	}

	details := s.formatter.FormatDetails(resources[0])
	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource described successfully")
	return details, nil
}
