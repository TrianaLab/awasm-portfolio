package service

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/util"
	"fmt"
)

type DeleteService struct {
	repo *repository.InMemoryRepository
}

func NewDeleteService(repo *repository.InMemoryRepository) *DeleteService {
	return &DeleteService{repo: repo}
}

func (s *DeleteService) DeleteResource(kind string, name string, namespace string) (string, error) {
	nKind, _ := util.NormalizeKind(kind)
	if nKind == "" {
		return "", fmt.Errorf("you must specify only one resource")
	}

	if name == "" {
		return "", fmt.Errorf("resource(s) were provided, but no name was specified")
	}

	if namespace == "" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

	if nKind == "namespace" {
		return s.repo.Delete("all", "", name)
	}

	return s.repo.Delete(kind, name, namespace)
}
