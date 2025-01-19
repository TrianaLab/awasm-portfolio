package service

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/util"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type DeleteService struct {
	repo *repository.InMemoryRepository
	cmd  *cobra.Command
}

func NewDeleteService(repo *repository.InMemoryRepository, cmd *cobra.Command) *DeleteService {
	return &DeleteService{
		repo: repo,
		cmd:  cmd,
	}
}

func (s *DeleteService) DeleteResource(kind, name, namespace string) (string, error) {
	nKind, err := util.NormalizeKind(kind)
	if err != nil {
		return "", err
	}
	if nKind == "" {
		return "", fmt.Errorf("you must specify only one resource")
	}

	if name == "" {
		return "", fmt.Errorf("resource(s) were provided, but no name was specified")
	}

	if namespace == "" && kind != "namespace" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

	var deletedResources []string

	if nKind == "namespace" {
		// First we delete the namespace
		deletedNamespace, err := s.repo.Delete("namespace", name, namespace)
		if err != nil {
			return "", err
		}

		// Then the inner resources
		deletedChildren, err := s.repo.Delete("all", "", name)
		if err != nil {
			return "", err
		}

		deletedResources = append(deletedResources, deletedNamespace, deletedChildren)

		return strings.Join(deletedResources, "\n"), nil
	}

	resources, err := s.repo.List("all", "", "")
	if err != nil {
		return "", err
	}

	r, err := s.repo.Delete(kind, name, namespace)
	if err != nil {
		return "", err
	}

	deletedResources = append(deletedResources, r)
	for _, res := range resources {
		if res.GetOwnerReference().GetID() == fmt.Sprintf("%s:%s:%s", nKind, name, namespace) {
			deleted, err := s.repo.Delete(res.GetKind(), res.GetName(), res.GetNamespace())
			if err != nil {
				return "", err
			}
			deletedResources = append(deletedResources, deleted)
		}
	}

	return strings.Join(deletedResources, "\n"), nil
}
