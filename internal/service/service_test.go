package service_test

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/service"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

type dummyResource struct {
	Kind              string    `yaml:"Kind,omitempty" json:"Kind,omitempty"`
	Name              string    `yaml:"Name,omitempty" json:"Name,omitempty"`
	Namespace         string    `yaml:"Namespace,omitempty" json:"Namespace,omitempty"`
	CreationTimestamp time.Time `yaml:"CreationTimestamp,omitempty" json:"CreationTimestamp,omitempty"`
	ownerRef          models.OwnerReference
}

func (d *dummyResource) GetKind() string                             { return d.Kind }
func (d *dummyResource) GetName() string                             { return d.Name }
func (d *dummyResource) SetName(name string)                         { d.Name = name }
func (d *dummyResource) GetNamespace() string                        { return d.Namespace }
func (d *dummyResource) SetNamespace(namespace string)               { d.Namespace = namespace }
func (d *dummyResource) GetOwnerReference() models.OwnerReference    { return d.ownerRef }
func (d *dummyResource) SetOwnerReference(ref models.OwnerReference) { d.ownerRef = ref }
func (d *dummyResource) GetID() string                               { return d.Kind + ":" + d.Name + ":" + d.Namespace }
func (d *dummyResource) GetCreationTimestamp() time.Time             { return d.CreationTimestamp }
func (d *dummyResource) SetCreationTimestamp(t time.Time)            { d.CreationTimestamp = t }

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("output", "table", "output format")
	return cmd
}

func newTestResource(kind, name, namespace string) models.Resource {
	return &dummyResource{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}
}

func TestCreateService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	cs := service.NewCreateService(repo, cmd)

	_, err := cs.CreateResource("invalidKind", "testName", "testNS")
	if err == nil || !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error, got %v", err)
	}

	_, err = cs.CreateResource("resume", "testName", "nonexistentNS")
	if err == nil || !strings.Contains(err.Error(), "namespace 'nonexistentNS' not found") {
		t.Errorf("expected namespace not found error, got %v", err)
	}

	namespace := newTestResource("namespace", "test", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err := cs.CreateResource("resume", "testResume", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "resume/testResume created") {
		t.Errorf("unexpected success message: %s", msg)
	}

	t.Run("Create Resume With Nested Resources", func(t *testing.T) {
		namespace := newTestResource("namespace", "nested-test", "")
		_, err := repo.Create(namespace)
		if err != nil {
			t.Fatal(err)
		}

		_, err = cs.CreateResource("resume", "test-nested", "nested-test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		types := []string{
			"work", "education", "volunteer", "award", "skill",
			"language", "interest", "reference", "project", "basics",
		}

		for _, typ := range types {
			resources, err := repo.List(typ, "", "nested-test")
			if err != nil {
				t.Errorf("error listing %s: %v", typ, err)
			}
			if len(resources) == 0 {
				t.Errorf("no %s resources created", typ)
			}

			for _, res := range resources {
				ownerRef := res.GetOwnerReference()
				if ownerRef.Kind != "resume" {
					t.Errorf("expected owner kind resume, got %s", ownerRef.Kind)
				}
				if ownerRef.Name != "test-nested" {
					t.Errorf("expected owner name test-nested, got %s", ownerRef.Name)
				}
				if ownerRef.Namespace != "nested-test" {
					t.Errorf("expected owner namespace nested-test, got %s", ownerRef.Namespace)
				}
				if res.GetNamespace() != "nested-test" {
					t.Errorf("expected namespace nested-test, got %s", res.GetNamespace())
				}
			}
		}
	})
}

func TestDeleteService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	ds := service.NewDeleteService(repo, cmd)

	_, err := ds.DeleteResource("invalidKind", "testName", "test")
	if err == nil || !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error, got %v", err)
	}

	_, err = ds.DeleteResource("resume", "", "test")
	if err == nil || !strings.Contains(err.Error(), "no name was specified") {
		t.Errorf("expected missing name error, got %v", err)
	}

	_, err = ds.DeleteResource("resume", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	namespace := newTestResource("namespace", "test", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resume := newTestResource("resume", "testResume", "test")
	_, err = repo.Create(resume)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err := ds.DeleteResource("namespace", "test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "namespace/test in namespace '' deleted") {
		t.Errorf("unexpected success message: %s", msg)
	}

	_, err = ds.DeleteResource("namespace", "inexistent", "")
	if err == nil || !strings.Contains(err.Error(), "namespace/inexistent not found in namespace ''") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	t.Run("Empty Kind", func(t *testing.T) {
		_, err := ds.DeleteResource("", "test", "default")
		if err == nil {
			t.Error("expected error for empty kind")
		}
		if !strings.Contains(err.Error(), "unsupported resource kind") {
			t.Errorf("expected 'unsupported resource kind' error, got %v", err)
		}
	})

	t.Run("Delete With Owner References", func(t *testing.T) {
		namespace := newTestResource("namespace", "test-ns", "")
		_, err := repo.Create(namespace)
		if err != nil {
			t.Fatal(err)
		}

		parent := newTestResource("resume", "parent-resume", "test-ns")
		_, err = repo.Create(parent)
		if err != nil {
			t.Fatal(err)
		}

		parentID := strings.ToLower("resume:parent-resume:test-ns")

		child1 := &dummyResource{
			Kind:      "work",
			Name:      "child-work",
			Namespace: "test-ns",
			ownerRef: models.OwnerReference{
				Kind:      "resume",
				Name:      "parent-resume",
				Namespace: "test-ns",
			},
		}
		_, err = repo.Create(child1)
		if err != nil {
			t.Fatal(err)
		}

		child2 := &dummyResource{
			Kind:      "education",
			Name:      "child-edu",
			Namespace: "test-ns",
			ownerRef: models.OwnerReference{
				Kind:      "resume",
				Name:      "parent-resume",
				Namespace: "test-ns",
			},
		}
		_, err = repo.Create(child2)
		if err != nil {
			t.Fatal(err)
		}

		resources, _ := repo.List("work", "", "test-ns")
		if len(resources) == 0 || resources[0].GetOwnerReference().GetID() != parentID {
			t.Fatal("child work not created with correct owner reference")
		}

		resources, _ = repo.List("education", "", "test-ns")
		if len(resources) == 0 || resources[0].GetOwnerReference().GetID() != parentID {
			t.Fatal("child education not created with correct owner reference")
		}

		msg, err := ds.DeleteResource("resume", "parent-resume", "test-ns")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(msg, "resume/parent-resume") {
			t.Error("parent resource not mentioned in delete message")
		}
		if !strings.Contains(msg, "work/child-work") {
			t.Error("child work not mentioned in delete message")
		}
		if !strings.Contains(msg, "education/child-edu") {
			t.Error("child education not mentioned in delete message")
		}

		resources, _ = repo.List("work", "", "test-ns")
		if len(resources) > 0 {
			t.Error("child work still exists after parent deletion")
		}

		resources, _ = repo.List("education", "", "test-ns")
		if len(resources) > 0 {
			t.Error("child education still exists after parent deletion")
		}
	})

	t.Run("No Resource Specified", func(t *testing.T) {
		_, err := ds.DeleteResource("all", "", "default")
		if err == nil {
			t.Error("expected error for no resource specified")
		}
		if !strings.Contains(err.Error(), "you must specify only one resource") {
			t.Errorf("expected 'you must specify only one resource' error, got %v", err)
		}
	})

	t.Run("Delete Error Cases", func(t *testing.T) {
		_, err := ds.DeleteResource("resume", "nonexistent", "nonexistent")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected 'not found' error, got %v", err)
		}

		namespace := newTestResource("namespace", "test-ns2", "")
		_, _ = repo.Create(namespace)

		parent := newTestResource("resume", "parent2", "test-ns2")
		_, _ = repo.Create(parent)

		child := &dummyResource{
			Kind:      "work",
			Name:      "child2",
			Namespace: "test-ns2",
			ownerRef: models.OwnerReference{
				Kind:      "resume",
				Name:      "parent2",
				Namespace: "test-ns2",
			},
		}
		_, _ = repo.Create(child)

		msg, err := ds.DeleteResource("work", "child2", "test-ns2")
		if err != nil {
			t.Errorf("unexpected error when deleting child: %v", err)
		}
		if !strings.Contains(msg, "work/child2") {
			t.Error("child deletion message not found")
		}

		msg, err = ds.DeleteResource("namespace", "test-ns2", "")
		if err != nil {
			t.Errorf("unexpected error when deleting namespace: %v", err)
		}
		if !strings.Contains(msg, "namespace/test-ns2") {
			t.Error("namespace deletion message not found")
		}
	})
}

func TestDescribeService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	ds := service.NewDescribeService(repo, cmd)

	_, err := ds.DescribeResource("resume", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	namespace := newTestResource("namespace", "testNS", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resume := newTestResource("resume", "testResume", "testNS")
	_, err = repo.Create(resume)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err2 := ds.DescribeResource("resume", "testResume", "testNS")
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if !strings.Contains(msg, "testResume") {
		t.Errorf("unexpected description output: %s", msg)
	}

	t.Run("Describe All Resources", func(t *testing.T) {
		resources := []struct {
			kind      string
			name      string
			namespace string
		}{
			{"work", "work-1", "testNS"},
			{"education", "edu-1", "testNS"},
			{"skill", "skill-1", "testNS"},
		}

		for _, r := range resources {
			res := newTestResource(r.kind, r.name, r.namespace)
			_, err := repo.Create(res)
			if err != nil {
				t.Fatalf("error creating resource %s/%s: %v", r.kind, r.name, err)
			}
		}

		for _, r := range resources {
			msg, err := ds.DescribeResource(r.kind, "", r.namespace)
			if err != nil {
				t.Errorf("error describing all %s: %v", r.kind, err)
			}
			if !strings.Contains(msg, r.name) {
				t.Errorf("expected %s in output, got: %s", r.name, msg)
			}
		}
	})

	t.Run("Describe All With Namespace", func(t *testing.T) {
		ns1 := newTestResource("namespace", "ns1", "")
		ns2 := newTestResource("namespace", "ns2", "")
		_, _ = repo.Create(ns1)
		_, _ = repo.Create(ns2)

		resources := []struct {
			kind      string
			name      string
			namespace string
		}{
			{"work", "work-1", "ns1"},
			{"work", "work-2", "ns2"},
			{"education", "edu-1", "ns1"},
			{"education", "edu-2", "ns2"},
		}

		for _, r := range resources {
			res := newTestResource(r.kind, r.name, r.namespace)
			_, _ = repo.Create(res)
		}

		msg, err := ds.DescribeResource("all", "", "ns1")
		if err != nil {
			t.Fatalf("error describing all resources: %v", err)
		}

		if strings.Contains(msg, "work-2") || strings.Contains(msg, "edu-2") {
			t.Error("found resources from wrong namespace")
		}
		if !strings.Contains(msg, "work-1") || !strings.Contains(msg, "edu-1") {
			t.Error("missing resources from correct namespace")
		}
		if !strings.Contains(msg, "ns1") || strings.Contains(msg, "ns2") {
			t.Error("incorrect namespace filtering")
		}
	})
}

func TestGetService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	err := cmd.Flags().Set("output", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gs := service.NewGetService(repo, cmd)

	_, err = gs.GetResources("resume", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	namespace := newTestResource("namespace", "testNS", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resume := newTestResource("resume", "testResume", "testNS")
	_, err = repo.Create(resume)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err := gs.GetResources("resume", "", "testNS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "testResume") {
		t.Errorf("unexpected retrieval output: %s", msg)
	}

	t.Run("Get Resources Without Namespaces", func(t *testing.T) {
		resources := []struct {
			kind      string
			name      string
			namespace string
		}{
			{"namespace", "ns1", ""},
			{"namespace", "ns2", ""},
			{"work", "work-1", "ns1"},
			{"education", "edu-1", "ns1"},
		}

		for _, r := range resources {
			res := newTestResource(r.kind, r.name, r.namespace)
			_, err := repo.Create(res)
			if err != nil {
				t.Fatalf("error creating resource %s/%s: %v", r.kind, r.name, err)
			}
		}

		err := cmd.Flags().Set("output", "json")
		if err != nil {
			t.Fatalf("error setting flag: %v", err)
		}

		msg, err := gs.GetResources("all", "", "")
		if err != nil {
			t.Fatalf("error getting resources: %v", err)
		}

		if strings.Contains(msg, "\"Kind\": \"namespace\"") {
			t.Error("namespaces should not be included in output")
		}
		if !strings.Contains(msg, "\"Kind\": \"work\"") {
			t.Error("work resources should be included in output")
		}
		if !strings.Contains(msg, "\"Kind\": \"education\"") {
			t.Error("education resources should be included in output")
		}
	})
}

func TestResourceServiceImpl(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	rs := service.NewResourceService(repo, cmd)

	_, _ = rs.CreateResource("resume", "testResume", "testNS")
	_, _ = rs.DeleteResource("resume", "testResume", "testNS")
	_, _ = rs.GetResources("resume", "testResume", "testNS")
	_, _ = rs.DescribeResource("resume", "testResume", "testNS")
}
