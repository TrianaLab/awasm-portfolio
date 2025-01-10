package repository

import (
	"awasm-portfolio/internal/models"
	"errors"
	"sync"
)

type InMemoryRepository struct {
	mu        sync.RWMutex
	resources map[string]map[string]models.Resource
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{resources: make(map[string]map[string]models.Resource)}
}

func (r *InMemoryRepository) Create(kind string, resource models.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[kind]; !exists {
		r.resources[kind] = make(map[string]models.Resource)
	}
	r.resources[kind][resource.GetName()] = resource
	return nil
}

func (r *InMemoryRepository) Get(kind, name string) (models.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, exists := r.resources[kind]; !exists {
		return nil, errors.New("resource kind not found")
	}
	resource, exists := r.resources[kind][name]
	if !exists {
		return nil, errors.New("resource not found")
	}
	return resource, nil
}

func (r *InMemoryRepository) List(kind string) []models.Resource {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var resources []models.Resource
	if kindResources, exists := r.resources[kind]; exists {
		for _, res := range kindResources {
			resources = append(resources, res)
		}
	}
	return resources
}

func (r *InMemoryRepository) Delete(kind, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[kind]; !exists {
		return errors.New("resource kind not found")
	}
	delete(r.resources[kind], name)
	return nil
}
