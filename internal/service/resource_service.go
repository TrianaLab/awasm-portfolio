package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
	"reflect"
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
	return s.DeleteResourceInNamespace(kind, name, "")
}

// DeleteResourceInNamespace deletes a resource in a specific namespace
func (s *ResourceService) DeleteResourceInNamespace(kind, name, namespace string) (string, error) {
	// Normalize kind to lowercase
	kind = strings.ToLower(kind)

	// Special handling for namespaces
	if kind == "namespace" {
		return s.CascadeDeleteNamespace(name)
	}

	// Fetch resource to ensure it exists
	resource, err := s.repo.Get(kind, name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			fmt.Printf("Resource %s/%s already deleted, skipping.\n", kind, name)
			return "", nil
		}
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	// Validate namespace if applicable
	if namespace != "" && resource.GetNamespace() != namespace {
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	// Perform cascading delete for child resources
	var messages []string
	children := s.repo.FindByOwner(kind, name, resource.GetNamespace())
	for _, child := range children {
		childKind := strings.ToLower(reflect.TypeOf(child).Elem().Name())
		fmt.Printf("Deleting child %s/%s in namespace '%s'\n", childKind, child.GetName(), resource.GetNamespace())

		msg, err := s.DeleteResourceInNamespace(childKind, child.GetName(), resource.GetNamespace())
		if err != nil {
			fmt.Printf("Error deleting child %s/%s: %s\n", childKind, child.GetName(), err)
			continue
		}
		messages = append(messages, msg)
	}

	// Delete the resource itself
	err = s.repo.Delete(kind, name)
	if err != nil {
		return "", err
	}

	messages = append(messages, fmt.Sprintf("%s/%s deleted successfully from namespace '%s'.", kind, name, resource.GetNamespace()))
	return strings.Join(messages, "\n"), nil
}

// CascadeDeleteNamespace handles deletion of a namespace and all its resources
func (s *ResourceService) CascadeDeleteNamespace(namespace string) (string, error) {
	var messages []string

	// Iterate over all resources in the namespace
	for kind, resources := range s.repo.ListAll() {
		for _, resource := range resources {
			if resource.GetNamespace() == namespace {
				fmt.Printf("Deleting %s/%s in namespace '%s'\n", kind, resource.GetName(), namespace)
				msg, err := s.DeleteResource(kind, resource.GetName())
				if err != nil {
					fmt.Printf("Error deleting %s/%s: %s\n", kind, resource.GetName(), err)
					continue
				}
				messages = append(messages, msg)
			}
		}
	}

	// Finally, delete the namespace itself
	err := s.repo.Delete("namespace", namespace)
	if err != nil {
		return "", err
	}
	messages = append(messages, fmt.Sprintf("namespace/%s deleted successfully.", namespace))

	return strings.Join(messages, "\n"), nil
}

// ListResources lists resources of a given kind with optional namespace filtering
func (s *ResourceService) ListResources(kind, namespace string, allNamespaces bool) ([]models.Resource, error) {
	resources := s.repo.List(kind)

	var filtered []models.Resource
	for _, res := range resources {
		// Filter by namespace if allNamespaces is false
		if allNamespaces || namespace == "" || strings.EqualFold(res.GetNamespace(), namespace) {
			filtered = append(filtered, res)
		}
	}

	// Return the filtered resources
	return filtered, nil
}

// DescribeResource provides detailed information about a specific resource
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

// ListAllResources lists all resources with optional namespace filtering
func (s *ResourceService) ListAllResources(namespace string, allNamespaces bool) (map[string][]models.Resource, error) {
	resourcesByKind := make(map[string][]models.Resource)

	// Iterate over all supported resource kinds
	for kind := range SupportedResources() {
		// Skip "all" as it's not a real resource kind
		if kind == "all" {
			continue
		}

		// Attempt to list resources with namespace filtering
		resources, err := s.ListResources(kind, namespace, allNamespaces)
		if err != nil {
			continue // Skip kinds with no resources
		}

		if len(resources) > 0 {
			resourcesByKind[kind] = resources
		}
	}

	return resourcesByKind, nil
}
