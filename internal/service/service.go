// Package service is the thin business-logic layer between the cobra
// command handlers and the repository / UI formatter. The four exported
// functions (Create, Delete, Get, Describe) cover the full surface; no
// struct services, no DI, no reflection.
package service

import (
	"fmt"
	"strings"

	"github.com/TrianaLab/awasm-portfolio/internal/factory"
	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/ui"
	"github.com/TrianaLab/awasm-portfolio/internal/util"
)

// Create constructs a resource of the given kind and persists it.
// Namespaces must exist before any other resource can be created within them.
func Create(repo *repository.InMemoryRepository, kind, name, namespace string) (string, error) {
	kind, err := util.NormalizeKind(kind)
	if err != nil {
		return "", err
	}

	if kind != "namespace" {
		resources, err := repo.List("namespace", namespace, "")
		if err != nil && len(resources) == 0 {
			return "", fmt.Errorf("failed to create %s/%s: namespace '%s' not found", kind, name, namespace)
		}
	}

	resource, err := factory.New(kind, name, namespace)
	if err != nil {
		return "", err
	}
	msg, err := repo.Create(resource)
	if err != nil {
		return "", err
	}

	// Resumes aggregate child resources — persist each one so kubectl get
	// <child-kind> shows them and the resume view sees a real graph.
	if resume, ok := resource.(*types.Resume); ok {
		if err := persistResumeChildren(repo, resume); err != nil {
			return "", err
		}
	}
	return msg, nil
}

// persistResumeChildren writes Basics + every slice element to the
// repository. Each child carries its own Meta (populated by the
// factory) including the owner reference back to the resume.
func persistResumeChildren(repo *repository.InMemoryRepository, r *types.Resume) error {
	children := []models.Resource{&r.Basics}
	for i := range r.Work {
		children = append(children, &r.Work[i])
	}
	for i := range r.Volunteer {
		children = append(children, &r.Volunteer[i])
	}
	for i := range r.Education {
		children = append(children, &r.Education[i])
	}
	for i := range r.Awards {
		children = append(children, &r.Awards[i])
	}
	for i := range r.Certificates {
		children = append(children, &r.Certificates[i])
	}
	for i := range r.Publications {
		children = append(children, &r.Publications[i])
	}
	for i := range r.Skills {
		children = append(children, &r.Skills[i])
	}
	for i := range r.Languages {
		children = append(children, &r.Languages[i])
	}
	for i := range r.Interests {
		children = append(children, &r.Interests[i])
	}
	for i := range r.References {
		children = append(children, &r.References[i])
	}
	for i := range r.Projects {
		children = append(children, &r.Projects[i])
	}
	for _, c := range children {
		if _, err := repo.Create(c); err != nil {
			return fmt.Errorf("failed to save %s/%s: %w", c.GetKind(), c.GetName(), err)
		}
	}
	return nil
}

// Delete removes a resource (and, for namespaces, every resource within
// it; for parent resources, every child that owns a reference back).
func Delete(repo *repository.InMemoryRepository, kind, name, namespace string) (string, error) {
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

	if nKind == "namespace" {
		return deleteNamespace(repo, name, namespace)
	}
	return deleteWithChildren(repo, nKind, name, namespace)
}

// Get returns a table/JSON/YAML rendering of matching resources.
func Get(repo *repository.InMemoryRepository, kind, name, namespace, output string) (string, error) {
	if name != "" && namespace == "" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

	resources, err := repo.List(kind, name, namespace)
	if err != nil {
		return "", err
	}

	// "all" hides namespace objects from the list — they are infrastructure.
	if strings.EqualFold(kind, "all") {
		resources = withoutNamespaces(resources)
	}

	return ui.FormatTable(resources, strings.ToLower(output)), nil
}

// Describe returns a detailed YAML rendering of matching resources.
func Describe(repo *repository.InMemoryRepository, kind, name, namespace string) (string, error) {
	if name != "" && namespace == "" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

	resources, err := repo.List(kind, name, namespace)
	if err != nil {
		return "", err
	}

	if strings.EqualFold(kind, "all") && namespace != "" {
		resources = onlyInNamespace(resources, namespace)
	}

	return ui.FormatDetails(resources), nil
}

func deleteNamespace(repo *repository.InMemoryRepository, name, namespace string) (string, error) {
	deletedNamespace, err := repo.Delete("namespace", name, namespace)
	if err != nil {
		return "", err
	}
	// Cluster-scope wildcard delete cannot fail: kind is normalised to "" and
	// no name is supplied, so List never errors.
	deletedChildren, _ := repo.Delete("all", "", name)
	return strings.Join([]string{deletedNamespace, deletedChildren}, "\n"), nil
}

func deleteWithChildren(repo *repository.InMemoryRepository, kind, name, namespace string) (string, error) {
	// List every resource so we can find owner-referenced children. The
	// wildcard list cannot fail for the same reason as above.
	allResources, _ := repo.List("all", "", "")

	r, err := repo.Delete(kind, name, namespace)
	if err != nil {
		return "", err
	}

	parentID := strings.ToLower(kind + ":" + name + ":" + namespace)
	out := []string{r}
	for _, res := range allResources {
		if res.GetOwnerReference().GetID() != parentID {
			continue
		}
		// Deleting a resource we just listed cannot fail.
		deleted, _ := repo.Delete(res.GetKind(), res.GetName(), res.GetNamespace())
		out = append(out, deleted)
	}
	return strings.Join(out, "\n"), nil
}

func withoutNamespaces(in []models.Resource) []models.Resource {
	out := make([]models.Resource, 0, len(in))
	for _, r := range in {
		if r.GetKind() != "namespace" {
			out = append(out, r)
		}
	}
	return out
}

func onlyInNamespace(in []models.Resource, namespace string) []models.Resource {
	out := make([]models.Resource, 0, len(in))
	for _, r := range in {
		if r.GetKind() == "namespace" && r.GetName() != namespace {
			continue
		}
		out = append(out, r)
	}
	return out
}
