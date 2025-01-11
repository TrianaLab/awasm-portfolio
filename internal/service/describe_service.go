package service

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
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

func (s *DescribeService) DescribeResource(kind, name, namespace string) (string, error) {
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", err
	}

	if namespace != "" && resource.GetNamespace() != namespace {
		return "", fmt.Errorf("resource %s/%s not found in namespace %s", kind, name, namespace)
	}

	return s.formatter.FormatDetails(resource), nil
}
