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

// Returns a list of resources matching the kind, name, and namespace criteria.
// In case of e
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
	return fmt.Sprintf("%s/%s created successfully in namespace '%s'", kind, resource.GetName(), resource.GetNamespace()), nil
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
