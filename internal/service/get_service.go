package service

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"fmt"

	"github.com/sirupsen/logrus"
)

type GetService struct {
	repo *repository.InMemoryRepository
}

func NewGetService(repo *repository.InMemoryRepository) *GetService {
	return &GetService{repo: repo}
}

func (s *GetService) GetResources(kind string, name string, namespace string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}).Trace("GetService.GetResources called")

	kind = NormalizeResourceName(kind)
	if !IsValidResource(kind) {
		logrus.WithFields(logrus.Fields{
			"kind": kind,
		}).Error("Unsupported resource kind")
		return "", fmt.Errorf("unsupported resource kind: %s", kind)
	}

	if name == "" {
		resources, err := s.repo.List(kind)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"kind":  kind,
				"error": err,
			}).Error("Failed to list resources")
			return "", err
		}

		var filtered []models.Resource
		for _, res := range resources {
			if namespace == "" || res.GetNamespace() == namespace {
				filtered = append(filtered, res)
			}
		}

		if len(filtered) == 0 {
			logrus.WithFields(logrus.Fields{
				"kind":      kind,
				"namespace": namespace,
			}).Info("No resources found")
			return "No resources found.", nil
		}

		var result string
		for _, res := range filtered {
			result += fmt.Sprintf("- %s/%s\n", kind, res.GetName())
		}

		logrus.WithFields(logrus.Fields{
			"kind":      kind,
			"namespace": namespace,
			"count":     len(filtered),
		}).Info("Resources retrieved successfully")
		return result, nil
	}

	resource, err := s.repo.Get(kind, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
			"error":     err,
		}).Error("Resource not found")
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	if namespace != "" && resource.GetNamespace() != namespace {
		logrus.WithFields(logrus.Fields{
			"kind":      kind,
			"name":      name,
			"namespace": namespace,
		}).Error("Resource found in a different namespace")
		return "", fmt.Errorf("%s/%s not found in namespace '%s'", kind, name, namespace)
	}

	logrus.WithFields(logrus.Fields{
		"kind":      kind,
		"name":      name,
		"namespace": namespace,
	}).Info("Resource retrieved successfully")
	return fmt.Sprintf("Name: %s\nNamespace: %s\n", resource.GetName(), resource.GetNamespace()), nil
}
