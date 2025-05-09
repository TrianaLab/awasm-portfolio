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

	// If the resource is a resume, dynamically process nested fields
	if resume, ok := resource.(*types.Resume); ok {
		val := reflect.ValueOf(resume).Elem()
		typ := val.Type()

		if basicsField := val.FieldByName("Basics"); basicsField.IsValid() {
			if basics, ok := basicsField.Addr().Interface().(models.Resource); ok {
				basics.SetOwnerReference(models.OwnerReference{
					Kind:      "resume",
					Name:      name,
					Namespace: namespace,
				})
				basics.SetNamespace(namespace)
				basics.SetName(strings.ToLower(fmt.Sprintf("%s-basics", name)))

				if _, err := s.repo.Create(basics); err != nil {
					return "", fmt.Errorf("failed to save basics: %w", err)
				}
			}
		}

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			if fieldType.Name == "Basics" {
				continue
			}

			if field.Kind() == reflect.Slice {
				elemType := field.Type().Elem()
				if reflect.PointerTo(elemType).Implements(reflect.TypeOf((*models.Resource)(nil)).Elem()) {
					for j := 0; j < field.Len(); j++ {
						elem := field.Index(j).Addr().Interface().(models.Resource)

						elem.SetOwnerReference(models.OwnerReference{
							Kind:      "resume",
							Name:      name,
							Namespace: namespace,
						})
						elem.SetNamespace(namespace)
						elem.SetName(strings.ToLower(fmt.Sprintf("%s-%s-%d", name, fieldType.Name, j)))

						if _, err := s.repo.Create(elem); err != nil {
							return "", fmt.Errorf("failed to save %s[%d]: %w", fieldType.Name, j, err)
						}
					}
				}
			}
		}
	}

	return msg, nil
}
