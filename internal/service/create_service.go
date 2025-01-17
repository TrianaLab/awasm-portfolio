package service

import (
	"awasm-portfolio/internal/factory"
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/util"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

type CreateService struct {
	repo    *repository.InMemoryRepository
	factory *factory.ResourceFactory
	cmd     *cobra.Command
}

func NewCreateService(repo *repository.InMemoryRepository, cmd *cobra.Command) *CreateService {
	return &CreateService{
		repo:    repo,
		factory: factory.NewResourceFactory(),
		cmd:     cmd,
	}
}

func (s *CreateService) CreateResource(kind string, name string, namespace string) (string, error) {
	// Normalize the kind
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		return "", err
	}

	// Ensure the namespace exists if the kind is not "namespace"
	if kind != "namespace" {
		resources, err := s.repo.List("namespace", namespace, "")
		if err != nil && len(resources) == 0 {
			return "", fmt.Errorf("failed to create %s/%s: namespace '%s' not found", kind, name, namespace)
		}
	}

	// Create the main resource using the factory
	resource := s.factory.Create(kind, map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	})

	// Save the main resource
	msg, err := s.repo.Create(resource)
	if err != nil {
		return "", err
	}

	// If the resource is a profile, dynamically process nested fields
	if profile, ok := resource.(*types.Profile); ok {
		// Use reflection to iterate over fields
		val := reflect.ValueOf(profile).Elem()
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			// Check if the field implements the `models.Resource` interface
			if field.Kind() == reflect.Struct && reflect.PointerTo(field.Type()).Implements(reflect.TypeOf((*models.Resource)(nil)).Elem()) {
				nested := field.Addr().Interface().(models.Resource)
				nested.SetOwnerReference(models.OwnerReference{
					Kind:      "profile",
					Name:      name,
					Namespace: namespace,
				})
				nested.SetNamespace(namespace)
				nested.SetName(strings.ToLower(fmt.Sprintf("%s-%s", name, fieldType.Name)))

				if _, err := s.repo.Create(nested); err != nil {
					return "", fmt.Errorf("failed to save %s: %w", fieldType.Name, err)
				}
			}
		}
	}

	return msg, nil
}
