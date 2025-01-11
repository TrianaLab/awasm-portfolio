package repository

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/util"
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

// List retrieves resources matching the kind, name, and namespace criteria.
func (r *InMemoryRepository) List(kind, name, namespace string) ([]models.Resource, error) {
	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "InMemoryRepository.List called")

	r.mu.RLock()
	defer r.mu.RUnlock()

	var resources []models.Resource
	for _, res := range r.resources {
		if matchResource(res, kind, name, namespace) {
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
	kind, err := util.NormalizeKind(resource.GetKind())
	if err != nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind in Create")
		return "", err
	}

	resourceID := resource.GetID()

	logger.Trace(logrus.Fields{
		"id":        resourceID,
		"kind":      kind,
		"name":      resource.GetName(),
		"namespace": resource.GetNamespace(),
	}, "InMemoryRepository.Create called")

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resourceID]; exists {
		err := fmt.Errorf("resource %s already exists", resourceID)
		logger.Error(logrus.Fields{
			"id":    resourceID,
			"error": err,
		}, "Failed to create resource")
		return "", err
	}

	if kind != "namespace" && resource.GetOwnerReference().Kind == "" {
		resource.SetOwnerReference(models.OwnerReference{
			Kind: "namespace",
			Name: resource.GetNamespace(),
		})
	}

	r.resources[resourceID] = resource
	logger.Info(logrus.Fields{
		"id": resourceID,
	}, "Resource created successfully")
	return fmt.Sprintf("%s/%s created successfully in namespace '%s'.", kind, resource.GetName(), resource.GetNamespace()), nil
}

func (r *InMemoryRepository) Delete(kind, name, namespace string) (string, error) {
	// Normalize and validate the kind
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind in Delete")
		return "", err
	}

	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}, "InMemoryRepository.Delete called")

	// Step 1: List resources to delete
	r.mu.RLock()
	resourcesToDelete, err := r.List(kind, name, namespace)
	r.mu.RUnlock()
	if err != nil || len(resourcesToDelete) == 0 {
		return "", fmt.Errorf("no resources found to delete")
	}

	// Step 2: Collect IDs and handle cascade deletions
	deletedResources := []string{}
	for _, res := range resourcesToDelete {
		resourceID := res.GetID()

		// Delete the resource
		r.mu.Lock()
		delete(r.resources, resourceID)
		r.mu.Unlock()

		deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", res.GetKind(), res.GetName(), res.GetNamespace()))
		logger.Info(logrus.Fields{
			"id": resourceID,
		}, "Resource deleted successfully")

		// Step 3: Handle cascade deletions
		if kind == "namespace" {
			namespaceToDelete := name

			// List all resources in the namespace and delete them
			cascadeResources, err := r.List("", "", namespaceToDelete) // Wildcard for kind
			if err != nil {
				logger.Error(logrus.Fields{
					"namespace": namespaceToDelete,
					"error":     err,
				}, "Failed to list resources for namespace cascade deletion")
				continue
			}

			for _, cascadeRes := range cascadeResources {
				cascadeID := cascadeRes.GetID()

				r.mu.Lock()
				delete(r.resources, cascadeID)
				r.mu.Unlock()

				deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", cascadeRes.GetKind(), cascadeRes.GetName(), cascadeRes.GetNamespace()))
				logger.Info(logrus.Fields{
					"id": cascadeID,
				}, "Cascade deleted resource in namespace")
			}
		} else {
			// For non-namespace resources, delete resources with matching OwnerReference
			cascadeResources, err := r.List("", "", "") // Wildcard for kind, name, and namespace
			if err != nil {
				logger.Error(logrus.Fields{
					"resourceID": resourceID,
					"error":      err,
				}, "Failed to list resources for owner reference cascade deletion")
				continue
			}

			for _, cascadeRes := range cascadeResources {
				ownerRef := cascadeRes.GetOwnerReference()
				if strings.ToLower(ownerRef.Kind) == kind &&
					strings.ToLower(ownerRef.Name) == strings.ToLower(name) &&
					strings.ToLower(ownerRef.Namespace) == strings.ToLower(namespace) {
					cascadeID := cascadeRes.GetID()

					r.mu.Lock()
					delete(r.resources, cascadeID)
					r.mu.Unlock()

					deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", cascadeRes.GetKind(), cascadeRes.GetName(), cascadeRes.GetNamespace()))
					logger.Info(logrus.Fields{
						"id": cascadeID,
					}, "Cascade deleted resource with matching owner reference")
				}
			}
		}
	}

	return fmt.Sprintf("Deleted resources:\n%s", strings.Join(deletedResources, "\n")), nil
}

// matchResource checks if a resource matches the given kind, name, and namespace criteria.
func matchResource(res models.Resource, kind, name, namespace string) bool {
	// Allow wildcard for kind by checking if kind is empty
	matchKind := kind == "" || strings.ToLower(res.GetKind()) == kind
	matchName := name == "" || strings.ToLower(res.GetName()) == name
	matchNamespace := namespace == "" || strings.ToLower(res.GetNamespace()) == namespace ||
		(res.GetNamespace() == "" && namespace != "")

	logger.Trace(logrus.Fields{
		"id":             res.GetID(),
		"kind":           res.GetKind(),
		"name":           res.GetName(),
		"namespace":      res.GetNamespace(),
		"matchKind":      matchKind,
		"matchName":      matchName,
		"matchNamespace": matchNamespace,
	}, "Matching resource")

	return matchKind && matchName && matchNamespace
}
