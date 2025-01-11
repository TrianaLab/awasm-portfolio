package repository

import (
	"awasm-portfolio/internal/models"
	"errors"
	"fmt"
	"sync"
)

type InMemoryRepository struct {
	mu        sync.RWMutex
	resources map[string]map[string]models.Resource
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		resources: make(map[string]map[string]models.Resource),
	}
}

// Create adds a new resource to the repository.
func (r *InMemoryRepository) Create(resource models.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	kind := resource.GetKind()
	name := resource.GetName()
	if _, exists := r.resources[kind]; !exists {
		r.resources[kind] = make(map[string]models.Resource)
	}

	if _, exists := r.resources[kind][name]; exists {
		return fmt.Errorf("%s/%s already exists", kind, name)
	}

	r.resources[kind][name] = resource
	return nil
}

// Get retrieves a single resource by kind and name.
func (r *InMemoryRepository) Get(kind, name string) (models.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		return nil, errors.New("resource kind not found")
	}

	resource, exists := resourcesByKind[name]
	if !exists {
		return nil, errors.New("resource not found")
	}

	return resource, nil
}

// Update modifies an existing resource.
func (r *InMemoryRepository) Update(resource models.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	kind := resource.GetKind()
	name := resource.GetName()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		return errors.New("resource kind not found")
	}

	if _, exists := resourcesByKind[name]; !exists {
		return errors.New("resource not found")
	}

	r.resources[kind][name] = resource
	return nil
}

// Delete removes a resource by kind and name.
func (r *InMemoryRepository) Delete(kind, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		return errors.New("resource kind not found")
	}

	if _, exists := resourcesByKind[name]; !exists {
		return fmt.Errorf("%s/%s not found", kind, name)
	}

	delete(resourcesByKind, name)
	return nil
}

// List retrieves all resources of a specific kind.
func (r *InMemoryRepository) List(kind string) ([]models.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		return nil, errors.New("resource kind not found")
	}

	var resources []models.Resource
	for _, resource := range resourcesByKind {
		resources = append(resources, resource)
	}

	return resources, nil
}

// ListAll retrieves all resources grouped by kind.
func (r *InMemoryRepository) ListAll() map[string][]models.Resource {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allResources := make(map[string][]models.Resource)
	for kind, resourcesByKind := range r.resources {
		var resources []models.Resource
		for _, resource := range resourcesByKind {
			resources = append(resources, resource)
		}
		allResources[kind] = resources
	}

	return allResources
}
