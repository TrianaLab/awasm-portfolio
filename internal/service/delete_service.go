package service

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/repository"
	"fmt"

	"github.com/sirupsen/logrus"
)

type DeleteService struct {
	repo *repository.InMemoryRepository
}

func NewDeleteService(repo *repository.InMemoryRepository) *DeleteService {
	return &DeleteService{repo: repo}
}

func (s *DeleteService) DeleteResource(kind string, name string, namespace string) (string, error) {
	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "DeleteService.DeleteResource called")

	if namespace == "" && kind != "namespace" {
		logger.Error(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "Namespace is required")
		return "", fmt.Errorf("namespace is required")
	}

	msg, err := s.repo.Delete(kind, name, namespace)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Failed to delete resource")
		return "", err
	}

	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource deleted successfully")
	return msg, nil
}
