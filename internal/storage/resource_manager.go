package storage

import (
	"awasm-portfolio/internal/models"
	"errors"
	"fmt"
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

	// Check for namespace existence if resource is namespaced
	if resource.Namespaced && resource.Namespace != "" {
		nsMap, exists := rm.Resources["namespace"]
		if !exists || nsMap[resource.Namespace].Name == "" {
			return fmt.Errorf("namespace '%s' does not exist", resource.Namespace)
		}
	}

	// Create resource type map if not exists
	if _, exists := rm.Resources[resourceType]; !exists {
		rm.Resources[resourceType] = make(map[string]models.ResourceBase)
	}

	// Avoid duplicates
	if _, exists := rm.Resources[resourceType][resource.Name]; exists {
		return fmt.Errorf("resource '%s' already exists", resource.Name)
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
