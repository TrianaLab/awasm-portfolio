package service

import (
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
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}).Trace("DeleteService.DeleteResource called")

	if namespace == "" {
		logrus.Error("Namespace is required")
		return "", fmt.Errorf("namespace is required")
	}

	err := s.repo.Delete(kind, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Failed to delete resource")
		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource deleted successfully")
	return fmt.Sprintf("%s/%s deleted successfully", kind, name), nil
}
