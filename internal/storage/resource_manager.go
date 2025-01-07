package storage

import (
	"awasm-portfolio/internal/models"
	"errors"
	"sync"
)

type ResourceManager struct {
	Resources map[string]map[string]models.ResourceBase
	mu        sync.RWMutex
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		Resources: make(map[string]map[string]models.ResourceBase),
	}
}

func (rm *ResourceManager) IsNamespaced(resourceType string) bool {
	namespacedResources := map[string]bool{
		"configmap": true,
		"secret":    true,
	}
	return namespacedResources[resourceType]
}

func (rm *ResourceManager) GetAll(resourceType, namespace string, allNamespaces bool) (map[string]models.ResourceBase, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	resourceMap, ok := rm.Resources[resourceType]
	if !ok {
		return nil, errors.New("resource type not found")
	}

	result := make(map[string]models.ResourceBase)
	for name, res := range resourceMap {
		if allNamespaces || !res.Namespaced || res.Namespace == namespace {
			result[name] = res
		}
	}

	return result, nil
}

func (rm *ResourceManager) Create(resourceType string, resource models.ResourceBase) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.Resources[resourceType]; !exists {
		rm.Resources[resourceType] = make(map[string]models.ResourceBase)
	}

	if _, exists := rm.Resources[resourceType][resource.Name]; exists {
		return errors.New("resource already exists")
	}

	rm.Resources[resourceType][resource.Name] = resource
	return nil
}

func (rm *ResourceManager) Delete(resourceType, name, namespace string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	resourceMap, ok := rm.Resources[resourceType]
	if !ok {
		return errors.New("resource type not found")
	}

	resource, exists := resourceMap[name]
	if !exists || (resource.Namespaced && resource.Namespace != namespace) {
		return errors.New("resource not found")
	}

	delete(resourceMap, name)
	return nil
}
