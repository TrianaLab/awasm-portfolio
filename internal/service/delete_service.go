package service

import (
	"fmt"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/util"
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

		// Then the inner resources. The "all" call cannot fail because the
		// kind normalizes to the empty wildcard and no name is supplied.
		deletedChildren, _ := s.repo.Delete("all", "", name)

		deletedResources = append(deletedResources, deletedNamespace, deletedChildren)

		return strings.Join(deletedResources, "\n"), nil
	}

	// The "all" wildcard list cannot fail; ignore the error.
	resources, _ := s.repo.List("all", "", "")

	r, err := s.repo.Delete(kind, name, namespace)
	if err != nil {
		return "", err
	}

	deletedResources = append(deletedResources, r)
	for _, res := range resources {
		if res.GetOwnerReference().GetID() == strings.ToLower(nKind+":"+name+":"+namespace) {
			// Deleting a resource we just listed cannot fail.
			deleted, _ := s.repo.Delete(res.GetKind(), res.GetName(), res.GetNamespace())
			deletedResources = append(deletedResources, deleted)
		}
	}

	return strings.Join(deletedResources, "\n"), nil
}
