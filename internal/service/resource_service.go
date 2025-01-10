package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
	"strings"
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

// CreateResource handles resource creation with duplicate validation
func (s *ResourceService) CreateResource(kind string, resource models.Resource) (string, error) {
	// Validate namespace existence
	if resource.GetNamespace() != "" && kind != "namespace" {
		_, err := s.repo.Get("namespace", resource.GetNamespace())
		if err != nil {
			return "", fmt.Errorf("namespace '%s' does not exist", resource.GetNamespace())
		}
	}

	// Check for duplicates
	existing, _ := s.repo.Get(kind, resource.GetName())
	if existing != nil && existing.GetNamespace() == resource.GetNamespace() {
		return "", fmt.Errorf("%s/%s already exists in namespace '%s'", kind, resource.GetName(), resource.GetNamespace())
	}

	// Proceed with resource creation
	err := s.repo.Create(kind, resource)
	if err != nil {
		return "", err
	}

	// Return creation message
	return fmt.Sprintf("%s/%s created successfully in namespace '%s'.\n", kind, resource.GetName(), resource.GetNamespace()), nil
}

// DeleteResource handles cascading deletion
func (s *ResourceService) DeleteResource(kind, name string) (string, error) {
	// Fetch resource to ensure it exists
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", err
	}

	var messages []string

	// Perform cascading delete for child resources
	for _, owner := range resource.GetOwnerReferences() {
		msg, err := s.DeleteResource(owner.Kind, owner.Name)
		if err != nil {
			return "", err
		}
		messages = append(messages, msg)
	}

	// Delete the resource itself
	err = s.repo.Delete(kind, name)
	if err != nil {
		return "", err
	}

	// Append deletion message
	messages = append(messages, fmt.Sprintf("%s/%s deleted successfully from namespace '%s'.", kind, name, resource.GetNamespace()))

	// Combine all messages
	return s.formatter.FormatDetails(resource) + "\n" + fmt.Sprintf(strings.Join(messages, "\n")), nil
}

// DeleteResourceWithNamespace validates namespace before deletion
func (s *ResourceService) DeleteResourceWithNamespace(kind, name, namespace string) (string, error) {
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		return "", err
	}

	// Validate namespace
	if namespace != "" && resource.GetNamespace() != namespace {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	return s.DeleteResource(kind, name)
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
