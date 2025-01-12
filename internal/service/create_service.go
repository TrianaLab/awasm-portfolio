package service

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/logger"
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
	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "CreateService.CreateResource called")

	if kind != "namespace" && !s.repo.Exists("namespace", namespace, "") {
		logger.Error(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "Namespace does not exist")
		return "", fmt.Errorf("namespace '%s' does not exist", namespace)
	}

	resource := s.factory.Create(kind, map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	})
	if resource == nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind")
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	msg, err := s.repo.Create(resource)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Failed to create resource")
		return "", err
	}

	logger.Info(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "Resource created successfully")
	return msg, nil
}
