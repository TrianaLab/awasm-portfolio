package repository

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/util"
	"fmt"
	"strings"
	"sync"
	"time"

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

	if len(resources) == 0 {
		err := fmt.Errorf("%s/%s was not found in '%s' namespace", kind, name, namespace)
		logger.Error(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}, "No resources found")
		return nil, err
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

	// Set the creation timestamp here
	if resource.GetCreationTimestamp().IsZero() {
		resource.SetCreationTimestamp(time.Now())
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

	var deletedResources []string

	// If the resource to be deleted is a namespace
	if kind == "namespace" {
		r.mu.RLock()
		allResources, err := r.List("", "", name) // Get all objects in the namespace
		r.mu.RUnlock()
		if err != nil {
			return "", fmt.Errorf("failed to list resources in namespace '%s': %v", name, err)
		}

		// Filter out namespaces except the one being deleted
		resourcesToDelete := []models.Resource{}
		for _, res := range allResources {
			resKind, _ := util.NormalizeKind(res.GetKind())
			if resKind != "namespace" || res.GetName() == name {
				resourcesToDelete = append(resourcesToDelete, res)
			}
		}

		// Delete the filtered resources
		for _, res := range resourcesToDelete {
			resourceID := res.GetID() // Use unique ID for deletion

			r.mu.Lock()
			if _, exists := r.resources[resourceID]; exists {
				delete(r.resources, resourceID) // Correct deletion
				deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", res.GetKind(), res.GetName(), name))
			}
			r.mu.Unlock()

			logger.Info(logrus.Fields{
				"kind": res.GetKind(),
				"name": res.GetName(),
			}, "Resource deleted in namespace cleanup")
		}

		return fmt.Sprintf("Deleted resources in namespace '%s':\n%s", name, strings.Join(deletedResources, "\n")), nil
	}

	// Step 1: List resources to delete
	r.mu.RLock()
	resourcesToDelete, err := r.List(kind, name, namespace)
	r.mu.RUnlock()
	if err != nil || len(resourcesToDelete) == 0 {
		return "", fmt.Errorf("no resources found to delete")
	}

	// Step 2: Delete resources and cascade to owner references
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

		// Find and delete resources with this resource as their owner
		r.mu.RLock()
		ownedResources, err := r.List("", "", "")
		r.mu.RUnlock()
		if err != nil {
			logger.Error(logrus.Fields{
				"ownerID": resourceID,
				"error":   err,
			}, "Failed to list resources for cascade deletion")
			continue
		}

		for _, ownedRes := range ownedResources {
			if ownedRes.GetOwnerReference().GetID() == resourceID { // Match owner references
				r.mu.Lock()
				delete(r.resources, ownedRes.GetID())
				r.mu.Unlock()

				deletedResources = append(deletedResources, fmt.Sprintf("%s/%s in namespace '%s'", ownedRes.GetKind(), ownedRes.GetName(), ownedRes.GetNamespace()))
				logger.Info(logrus.Fields{
					"id": ownedRes.GetID(),
				}, "Cascade deleted resource")
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

func (r *InMemoryRepository) Exists(kind, name, namespace string) bool {
	// Normalize and validate the kind
	normalizedKind, err := util.NormalizeKind(kind)
	if err != nil {
		logger.Error(logrus.Fields{
			"kind": kind,
		}, "Unsupported resource kind in Exists")
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, res := range r.resources {
		matchKind := normalizedKind == "" || strings.ToLower(res.GetKind()) == normalizedKind
		matchName := name == "" || strings.ToLower(res.GetName()) == strings.ToLower(name)
		matchNamespace := namespace == "" || strings.ToLower(res.GetNamespace()) == strings.ToLower(namespace)

		if matchKind && matchName && matchNamespace {
			return true
		}
	}
	return false
}
