package service

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"

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

	resource, err := s.repo.Get(kind, name)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Failed to describe resource")
		return "", err
	}

	details := fmt.Sprintf("Name: %s\nNamespace: %s\nKind: %s\n", resource.GetName(), resource.GetNamespace(), kind)
	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource described successfully")
	return details, nil
}
