package service

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/repository"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CreateService struct {
	repo    *repository.InMemoryRepository
	factory *factory.ResourceFactory
}

func NewCreateService(repo *repository.InMemoryRepository) *CreateService {
	return &CreateService{
		repo:    repo,
		factory: factory.NewResourceFactory(),
	}
}

func (s *CreateService) CreateResource(kind string, name string, namespace string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}).Trace("CreateService.CreateResource called")

	if namespace == "" {
		logrus.Error("Namespace is required")
		return "", fmt.Errorf("namespace is required")
	}

	resource := s.factory.Create(kind, map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	})
	if resource == nil {
		logrus.WithFields(logrus.Fields{
			"kind": kind,
		}).Error("Unsupported resource kind")
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}
	err := s.repo.Create(resource)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Failed to create resource")
		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource created successfully")
	return fmt.Sprintf("%s/%s created successfully in namespace '%s'.", kind, name, namespace), nil
}
