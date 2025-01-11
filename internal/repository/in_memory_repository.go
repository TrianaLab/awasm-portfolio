package repository

import (
	"awasm-portfolio/internal/logger"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/util"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
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

func (r *InMemoryRepository) Create(resource models.Resource) error {
	kind := util.NormalizeResourceName(resource.GetKind()) // Normalize kind
	name := resource.GetName()

	logger.Trace(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": resource.GetNamespace(),
	}, "InMemoryRepository.Create called")

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[kind]; !exists {
		r.resources[kind] = make(map[string]models.Resource)
	}
	if _, exists := r.resources[kind][name]; exists {
		err := fmt.Errorf("%s/%s already exists", kind, name)
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Resource already exists")
		return err
	}

	r.resources[kind][name] = resource
	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource created successfully")
	return nil
}

func (r *InMemoryRepository) Get(kind, name string) (models.Resource, error) {
	logger.Trace(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "InMemoryRepository.Get called")

	r.mu.RLock()
	defer r.mu.RUnlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		err := errors.New("resource kind not found")
		logger.Error(logrus.Fields{
			"kind":  kind,
			"error": err,
		}, "Resource kind not found")
		return nil, err
	}

	resource, exists := resourcesByKind[name]
	if !exists {
		err := fmt.Errorf("resource %s/%s not found", kind, name)
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Resource not found")
		return nil, err
	}

	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource retrieved successfully")
	return resource, nil
}

func (r *InMemoryRepository) Update(resource models.Resource) error {
	logger.Trace(logrus.Fields{
		"kind":      resource.GetKind(),
		"name":      resource.GetName(),
		"namespace": resource.GetNamespace(),
	}, "InMemoryRepository.Update called")

	r.mu.Lock()
	defer r.mu.Unlock()

	kind := resource.GetKind()
	name := resource.GetName()

	resourcesByKind, exists := r.resources[kind]
	if !exists || resourcesByKind[name] == nil {
		err := fmt.Errorf("resource %s/%s not found", kind, name)
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Resource not found")
		return err
	}

	r.resources[kind][name] = resource
	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource updated successfully")
	return nil
}

func (r *InMemoryRepository) Delete(kind, name string) error {
	logger.Trace(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "InMemoryRepository.Delete called")

	r.mu.Lock()
	defer r.mu.Unlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists || resourcesByKind[name] == nil {
		err := fmt.Errorf("resource %s/%s not found", kind, name)
		logger.Error(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}, "Resource not found")
		return err
	}

	delete(resourcesByKind, name)
	logger.Info(logrus.Fields{
		"kind": kind,
		"name": name,
	}, "Resource deleted successfully")
	return nil
}

func (r *InMemoryRepository) List(kind string) ([]models.Resource, error) {
	logger.Trace(logrus.Fields{
		"kind": kind,
	}, "InMemoryRepository.List called")

	r.mu.RLock()
	defer r.mu.RUnlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		err := fmt.Errorf("resource kind %s not found", kind)
		logger.Error(logrus.Fields{
			"kind":  kind,
			"error": err,
		}, "Resource kind not found")
		return nil, err
	}

	var resources []models.Resource
	for _, res := range resourcesByKind {
		resources = append(resources, res)
	}

	logger.Info(logrus.Fields{
		"kind":  kind,
		"count": len(resources),
	}, "Resources listed successfully")
	return resources, nil
}
