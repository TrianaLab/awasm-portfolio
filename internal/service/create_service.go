package service

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/util"
	"fmt"
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
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		return "", err
	}

	if kind != "namespace" {
		resources, err := s.repo.List("namespace", namespace, "")
		if err != nil {
			return "", err
		}
		if len(resources) == 0 {
			return "", fmt.Errorf("failed to create %s: namespace '%s' not found", name, namespace)
		}
	}

	resource := s.factory.Create(kind, map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	})

	msg, err := s.repo.Create(resource)
	if err != nil {
		return "", err
	}

	return msg, nil
}
