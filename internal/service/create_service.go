package service

import (
	"fmt"
	"github.com/TrianaLab/awasm-portfolio/internal/factory"
	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/util"
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
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		return "", err
	}

	if kind != "namespace" {
		resources, err := s.repo.List("namespace", namespace, "")
		if err != nil && len(resources) == 0 {
			return "", fmt.Errorf("failed to create %s/%s: namespace '%s' not found", kind, name, namespace)
		}
	}

	resource := s.factory.Create(kind, map[string]interface{}{
		"name":      name,
		"namespace": namespace,
	})

	msg, err := s.repo.Create(resource)
	if err != nil {
		return "", err
	}

	if resume, ok := resource.(*types.Resume); ok {
		if err := s.createResumeChildren(resume, name, namespace); err != nil {
			return "", err
		}
	}

	return msg, nil
}

// createResumeChildren walks a Resume's fields via reflection and persists the
// derived child resources (Basics + each Resource-implementing slice element).
func (s *CreateService) createResumeChildren(resume *types.Resume, name, namespace string) error {
	val := reflect.ValueOf(resume).Elem()
	typ := val.Type()

	if basicsField := val.FieldByName("Basics"); basicsField.IsValid() {
		if basics, ok := basicsField.Addr().Interface().(models.Resource); ok {
			s.stampChild(basics, name, namespace, strings.ToLower(fmt.Sprintf("%s-basics", name)))
			if _, err := s.repo.Create(basics); err != nil {
				return fmt.Errorf("failed to save basics: %w", err)
			}
		}
	}

	// Every []T slice field on Resume holds Resource-implementing types by
	// construction; no need to re-check via reflection at runtime.
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		if fieldType.Name == "Basics" || field.Kind() != reflect.Slice {
			continue
		}
		if err := s.persistSlice(field, fieldType.Name, name, namespace); err != nil {
			return err
		}
	}
	return nil
}

// persistSlice iterates a reflected slice of Resource elements and saves each
// one with a derived name and owner reference back to the parent resume.
func (s *CreateService) persistSlice(field reflect.Value, fieldName, name, namespace string) error {
	for j := 0; j < field.Len(); j++ {
		elem := field.Index(j).Addr().Interface().(models.Resource)
		s.stampChild(elem, name, namespace, strings.ToLower(fmt.Sprintf("%s-%s-%d", name, fieldName, j)))
		if _, err := s.repo.Create(elem); err != nil {
			return fmt.Errorf("failed to save %s[%d]: %w", fieldName, j, err)
		}
	}
	return nil
}

// stampChild assigns owner reference, namespace, and derived name to a child resource.
func (s *CreateService) stampChild(child models.Resource, parentName, namespace, childName string) {
	child.SetOwnerReference(models.OwnerReference{
		Kind:      "resume",
		Name:      parentName,
		Namespace: namespace,
	})
	child.SetNamespace(namespace)
	child.SetName(childName)
}
