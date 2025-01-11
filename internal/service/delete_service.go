package service

import (
	"awasm-portfolio/internal/repository"
	"fmt"
)

type DeleteService struct {
	repo *repository.InMemoryRepository
}

func NewDeleteService(repo *repository.InMemoryRepository) *DeleteService {
	return &DeleteService{repo: repo}
}

func (s *DeleteService) DeleteResource(kind, name, namespace string) (string, error) {
	if namespace == "" {
		return "", fmt.Errorf("delete command does not support --all-namespaces")
	}

	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	if resource.GetNamespace() != namespace {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	// Delete the resource
	err = s.repo.Delete(kind, name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s deleted successfully from namespace '%s'.", kind, name, namespace), nil
}
