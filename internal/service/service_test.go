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

// Dummy implementation for models.Resource
type dummyResource struct {
	Kind              string    `yaml:"Kind,omitempty" json:"Kind,omitempty"`
	Name              string    `yaml:"Name,omitempty" json:"Name,omitempty"`
	Namespace         string    `yaml:"Namespace,omitempty" json:"Namespace,omitempty"`
	CreationTimestamp time.Time `yaml:"CreationTimestamp,omitempty" json:"CreationTimestamp,omitempty"`
}

func (d *dummyResource) GetKind() string                          { return d.Kind }
func (d *dummyResource) GetName() string                          { return d.Name }
func (d *dummyResource) SetName(name string)                      { d.Name = name }
func (d *dummyResource) GetNamespace() string                     { return d.Namespace }
func (d *dummyResource) SetNamespace(namespace string)            { d.Namespace = namespace }
func (d *dummyResource) GetOwnerReference() models.OwnerReference { return models.OwnerReference{} }
func (d *dummyResource) SetOwnerReference(models.OwnerReference)  {}
func (d *dummyResource) GetID() string                            { return d.Kind + ":" + d.Name + ":" + d.Namespace }
func (d *dummyResource) GetCreationTimestamp() time.Time          { return d.CreationTimestamp }
func (d *dummyResource) SetCreationTimestamp(t time.Time)         { d.CreationTimestamp = t }

// Helper to create a dummy command for testing
func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("output", "table", "output format")
	return cmd
}

// Helper to create a test resource
func newTestResource(kind, name, namespace string) models.Resource {
	return &dummyResource{
		Kind:      kind,
		Name:      name,
		Namespace: namespace,
	}
}

// Test for CreateService
func TestCreateService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	cs := service.NewCreateService(repo, cmd)

	// Test invalid kind
	_, err := cs.CreateResource("invalidKind", "testName", "testNS")
	if err == nil || !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error, got %v", err)
	}

	// Test namespace not found
	_, err = cs.CreateResource("resume", "testName", "nonexistentNS")
	if err == nil || !strings.Contains(err.Error(), "namespace 'nonexistentNS' not found") {
		t.Errorf("expected namespace not found error, got %v", err)
	}

	// Test successful creation
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
}

// Test for DeleteService
func TestDeleteService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	ds := service.NewDeleteService(repo, cmd)

	// Test invalid kind
	_, err := ds.DeleteResource("invalidKind", "testName", "test")
	if err == nil || !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("expected unsupported resource kind error, got %v", err)
	}

	// Test missing name
	_, err = ds.DeleteResource("resume", "", "test")
	if err == nil || !strings.Contains(err.Error(), "no name was specified") {
		t.Errorf("expected missing name error, got %v", err)
	}

	// Test missing namespace
	_, err = ds.DeleteResource("resume", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	// Test successful deletion
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

	// Test delete inexistent namespace
	_, err = ds.DeleteResource("namespace", "inexistent", "")
	if err == nil || !strings.Contains(err.Error(), "namespace/inexistent not found in namespace ''") {
		t.Errorf("expected missing namespace error, got %v", err)
	}
}

// Test for DescribeService
func TestDescribeService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	ds := service.NewDescribeService(repo, cmd)

	// Test missing namespace
	_, err := ds.DescribeResource("resume", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	// Test successful description
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

// Test for GetService
func TestGetService(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	err := cmd.Flags().Set("output", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gs := service.NewGetService(repo, cmd)

	// Test missing namespace
	_, err = gs.GetResources("resume", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	// Test successful retrieval
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
}

// Test for ResourceServiceImpl
func TestResourceServiceImpl(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	rs := service.NewResourceService(repo, cmd)

	// Test all delegations
	_, _ = rs.CreateResource("resume", "testResume", "testNS")
	_, _ = rs.DeleteResource("resume", "testResume", "testNS")
	_, _ = rs.GetResources("resume", "testResume", "testNS")
	_, _ = rs.DescribeResource("resume", "testResume", "testNS")
}
