package repository

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/util"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
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

// normalizeID normalizes the ID components to ensure consistent storage and retrieval.
func normalizeID(kind, name, namespace string) (string, string, string) {
	return strings.ToLower(kind), strings.ToLower(name), strings.ToLower(namespace)
}

// generateResourceID generates a unique ID for a resource based on its kind, name, and namespace.
func generateResourceID(kind, name, namespace string) string {
	kind, name, namespace = normalizeID(kind, name, namespace)
	return fmt.Sprintf("%s:%s:%s", kind, name, namespace)
}

// List retrieves resources matching the kind, name, and namespace criteria.
func (r *InMemoryRepository) List(kind, name, namespace string) ([]models.Resource, error) {
	kind = strings.ToLower(kind)
	name = strings.ToLower(name)
	namespace = strings.ToLower(namespace)

	if !util.IsValidResource(kind) {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind")
		return nil, fmt.Errorf("unsupported resource kind: %s", kind)
	}

	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "InMemoryRepository.List called")

	if kind == "" {
		logger.Error(logrus.Fields{"kind": kind}, "Kind is required")
		return nil, errors.New("kind is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var resources []models.Resource
	for id, res := range r.resources {
		matchKind := strings.ToLower(res.GetKind()) == kind
		matchName := name == "" || strings.ToLower(res.GetName()) == name
		matchNamespace := namespace == "" || res.GetNamespace() == namespace || (res.GetNamespace() == "" && namespace != "")

		logger.Trace(logrus.Fields{
			"id":             id,
			"matchKind":      matchKind,
			"matchName":      matchName,
			"matchNamespace": matchNamespace,
		}, "Matching resource")

		if matchKind && matchName && matchNamespace {
			resources = append(resources, res)
		}
	}

	logger.Info(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
		"count":     len(resources),
	}, "Resources listed successfully")
	return resources, nil
}

// Create adds a new resource to the repository.
func (r *InMemoryRepository) Create(resource models.Resource) (string, error) {
	kind := resource.GetKind()
	name := resource.GetName()
	namespace := resource.GetNamespace()

	if !util.IsValidResource(kind) {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind")
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	kind, name, namespace = normalizeID(kind, name, namespace)
	id := generateResourceID(kind, name, namespace)

	logger.Trace(logrus.Fields{
		"id":        id,
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "InMemoryRepository.Create called")

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[id]; exists {
		err := fmt.Errorf("resource %s already exists", id)
		logger.Error(logrus.Fields{
			"id":    id,
			"error": err,
		}, "Failed to create resource")
		return "", err
	}

	if kind != "namespace" && resource.GetOwnerReference().Kind == "" {
		resource.SetOwnerReference(models.OwnerReference{
			Kind: "namespace",
			Name: namespace,
		})
	}

	r.resources[id] = resource
	logger.Info(logrus.Fields{
		"id": id,
	}, "Resource created successfully")
	return fmt.Sprintf("%s/%s created successfully in namespace '%s'.", kind, name, namespace), nil
}

// Delete removes a resource and handles cascade deletions.
func (r *InMemoryRepository) Delete(kind, name, namespace string) (string, error) {
	kind = strings.ToLower(kind)
	name = strings.ToLower(name)
	namespace = strings.ToLower(namespace)

	if !util.IsValidResource(kind) {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind")
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "InMemoryRepository.Delete called")

	id := generateResourceID(kind, name, namespace)

	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.resources[id]
	if !exists {
		return "", fmt.Errorf("resource %s not found", id)
	}

	// Delete the resource
	delete(r.resources, id)
	logger.Info(logrus.Fields{
		"id": id,
	}, "Resource deleted successfully")

	deletedResources := []string{fmt.Sprintf("%s/%s in namespace '%s'", kind, name, namespace)}

	if kind == "namespace" {
		// Cascade delete all resources in the namespace
		for resID, res := range r.resources {
			if res.GetNamespace() == namespace {
				delete(r.resources, resID)
				deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", res.GetKind(), res.GetName(), res.GetNamespace()))
				logger.Info(logrus.Fields{
					"id": resID,
				}, "Cascade deleted resource in namespace")
			}
		}
	} else {
		// Cascade delete based on owner reference
		for resID, res := range r.resources {
			owner := res.GetOwnerReference()
			if owner.Kind == kind && owner.Name == name && owner.Namespace == namespace {
				delete(r.resources, resID)
				deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", res.GetKind(), res.GetName(), res.GetNamespace()))
				logger.Info(logrus.Fields{
					"id": resID,
				}, "Cascade deleted resource")
			}
		}
	}

	return fmt.Sprintf("Deleted resources:\n%s", stringList(deletedResources)), nil
}

// stringList formats a list of strings into a single string with newlines.
func stringList(items []string) string {
	return fmt.Sprintf("%s", strings.Join(items, "\n"))
}
