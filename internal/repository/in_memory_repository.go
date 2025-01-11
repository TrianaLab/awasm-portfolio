package repository

import (
	"awasm-portfolio/internal/models"
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
	logrus.WithFields(logrus.Fields{
		"kind":      resource.GetKind(),
		"name":      resource.GetName(),
		"namespace": resource.GetNamespace(),
	}).Trace("InMemoryRepository.Create called")

	r.mu.Lock()
	defer r.mu.Unlock()

	kind := resource.GetKind()
	name := resource.GetName()

	if _, exists := r.resources[kind]; !exists {
		r.resources[kind] = make(map[string]models.Resource)
	}
	if _, exists := r.resources[kind][name]; exists {
		err := fmt.Errorf("%s/%s already exists", kind, name)
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Resource already exists")
		return err
	}

	r.resources[kind][name] = resource
	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource created successfully")
	return nil
}

func (r *InMemoryRepository) Get(kind, name string) (models.Resource, error) {
	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Trace("InMemoryRepository.Get called")

	r.mu.RLock()
	defer r.mu.RUnlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		err := errors.New("resource kind not found")
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"error": err,
		}).Error("Resource kind not found")
		return nil, err
	}

	resource, exists := resourcesByKind[name]
	if !exists {
		err := fmt.Errorf("resource %s/%s not found", kind, name)
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Resource not found")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource retrieved successfully")
	return resource, nil
}

func (r *InMemoryRepository) Update(resource models.Resource) error {
	logrus.WithFields(logrus.Fields{
		"kind":      resource.GetKind(),
		"name":      resource.GetName(),
		"namespace": resource.GetNamespace(),
	}).Trace("InMemoryRepository.Update called")

	r.mu.Lock()
	defer r.mu.Unlock()

	kind := resource.GetKind()
	name := resource.GetName()

	resourcesByKind, exists := r.resources[kind]
	if !exists || resourcesByKind[name] == nil {
		err := fmt.Errorf("resource %s/%s not found", kind, name)
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Resource not found")
		return err
	}

	r.resources[kind][name] = resource
	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource updated successfully")
	return nil
}

func (r *InMemoryRepository) Delete(kind, name string) error {
	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Trace("InMemoryRepository.Delete called")

	r.mu.Lock()
	defer r.mu.Unlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists || resourcesByKind[name] == nil {
		err := fmt.Errorf("resource %s/%s not found", kind, name)
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"name":  name,
			"error": err,
		}).Error("Resource not found")
		return err
	}

	delete(resourcesByKind, name)
	logrus.WithFields(logrus.Fields{
		"kind": kind,
		"name": name,
	}).Info("Resource deleted successfully")
	return nil
}

func (r *InMemoryRepository) List(kind string) ([]models.Resource, error) {
	logrus.WithFields(logrus.Fields{
		"kind": kind,
	}).Trace("InMemoryRepository.List called")

	r.mu.RLock()
	defer r.mu.RUnlock()

	resourcesByKind, exists := r.resources[kind]
	if !exists {
		err := fmt.Errorf("resource kind %s not found", kind)
		logrus.WithFields(logrus.Fields{
			"kind":  kind,
			"error": err,
		}).Error("Resource kind not found")
		return nil, err
	}

	var resources []models.Resource
	for _, res := range resourcesByKind {
		resources = append(resources, res)
	}

	logrus.WithFields(logrus.Fields{
		"kind":  kind,
		"count": len(resources),
	}).Info("Resources listed successfully")
	return resources, nil
}
