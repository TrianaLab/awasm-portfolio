package service

import (
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
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}).Trace("DescribeService.DescribeResource called")

	resource, err := s.repo.Get(kind, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Failed to describe resource")
		return "", err
	}

	details := fmt.Sprintf("Name: %s\nNamespace: %s\nKind: %s\n", resource.GetName(), resource.GetNamespace(), kind)
	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource described successfully")
	return details, nil
}
