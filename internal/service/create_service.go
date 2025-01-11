package service

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/repository"
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

func (s *CreateService) CreateResource(kind, name, namespace string) (string, error) {
	if namespace == "" {
		return "", fmt.Errorf("create command does not support --all-namespaces")
	}

	// Validate namespace existence
	if kind != "namespace" {
		_, err := s.repo.Get("namespace", namespace)
		if err != nil {
			return "", fmt.Errorf("namespace '%s' does not exist", namespace)
		}
	}

	// Check for duplicates
	existing, _ := s.repo.Get(kind, name)
	if existing != nil && existing.GetNamespace() == namespace {
		return "", fmt.Errorf("%s/%s already exists in namespace '%s'", kind, name, namespace)
	}

	// Use the factory to create the resource
	resource := s.factory.Create(kind, map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	})
	if resource == nil {
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	// Save the resource to the repository
	err := s.repo.Create(resource)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s created successfully in namespace '%s'.", kind, name, namespace), nil
}
