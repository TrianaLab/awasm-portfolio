package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
)

type ResourceService struct {
	repo      *repository.InMemoryRepository
	formatter ui.Formatter
}

func NewResourceService(repo *repository.InMemoryRepository, formatter ui.Formatter) *ResourceService {
	return &ResourceService{
		repo:      repo,
		formatter: formatter,
	}
}

func (s *ResourceService) CreateResource(kind string, resource models.Resource) error {
	return s.repo.Create(kind, resource)
}

func (s *ResourceService) DeleteResource(kind, name string) error {
	// Fetch resource to ensure it exists
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return err
	}

	// Perform cascading delete for child resources
	for _, owner := range resource.GetOwnerReferences() {
		if err := s.DeleteResource(owner.Kind, owner.Name); err != nil {
			return err
		}
	}

	// Delete the resource itself
	return s.repo.Delete(kind, name)
}

func (s *ResourceService) ListResources(kind, namespace string, allNamespaces bool) ([]models.Resource, error) {
	resources, err := s.repo.List(kind)
	if err != nil {
		return nil, err
	}

	if allNamespaces {
		return resources, nil
	}

	// Filter by namespace
	var filtered []models.Resource
	for _, res := range resources {
		if res.GetNamespace() == namespace {
			filtered = append(filtered, res)
		}
	}

	return filtered, nil
}

func (s *ResourceService) DescribeResource(kind, name, namespace string) (string, error) {
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", err
	}

	// Ensure namespace matches
	if namespace != "" && resource.GetNamespace() != namespace {
		return "", fmt.Errorf("resource %s/%s not found in namespace %s", kind, name, namespace)
	}

	// Use the formatter to return resource details
	return s.formatter.FormatDetails(resource), nil
}

func (s *ResourceService) DeleteResourceWithNamespace(kind, name, namespace string) error {
	// Fetch resource to ensure it exists and matches the namespace
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return err
	}

	// Ensure namespace matches
	if namespace != "" && resource.GetNamespace() != namespace {
		return fmt.Errorf("%s/%s not found in namespace %s", kind, name, namespace)
	}

	// Perform cascading delete
	return s.DeleteResource(kind, name)
}
