package repository

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/util"
	"fmt"
	"strings"
	"sync"
	"time"
)

type InMemoryRepository struct {
	mu        sync.RWMutex
	resources map[string]models.Resource
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		resources: make(map[string]models.Resource),
	}
}

func (r *InMemoryRepository) List(kind, name, namespace string) ([]models.Resource, error) {
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var resources []models.Resource
	for _, res := range r.resources {
		if (kind == "" || strings.ToLower(res.GetKind()) == kind) &&
			(name == "" || strings.ToLower(res.GetName()) == name) &&
			(namespace == "" || strings.ToLower(res.GetNamespace()) == namespace || (res.GetNamespace() == "" && namespace != "")) {
			resources = append(resources, res)
		}
	}

	if len(resources) == 0 && name != "" {
		return nil, fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	return resources, nil
}

func (r *InMemoryRepository) Create(resource models.Resource) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resource.GetID()]; exists {
		return "", fmt.Errorf("failed to create %s: '%s' already exists", resource.GetKind(), resource.GetName())
	}

	if resource.GetCreationTimestamp().IsZero() {
		resource.SetCreationTimestamp(time.Now())
	}

	r.resources[resource.GetID()] = resource
	return fmt.Sprintf("%s/%s created", resource.GetKind(), resource.GetName()), nil
}

func (r *InMemoryRepository) Delete(kind, name, namespace string) (string, error) {
	r.mu.RLock()
	resources, err := r.List(kind, name, namespace)
	r.mu.RUnlock()
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	var deletedResources []string
	for _, res := range resources {
		if res.GetKind() != "namespace" || (res.GetKind() == "namespace" && res.GetName() == namespace) {
			delete(r.resources, res.GetID())
			deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s' deleted", res.GetKind(), res.GetName(), res.GetNamespace()))
		}
	}
	return fmt.Sprintf("%s", strings.Join(deletedResources, "\n")), nil
}
