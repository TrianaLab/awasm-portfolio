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
	_, err = cs.CreateResource("profile", "testName", "nonexistentNS")
	if err == nil || !strings.Contains(err.Error(), "namespace 'nonexistentNS' not found") {
		t.Errorf("expected namespace not found error, got %v", err)
	}

	// Test successful creation
	namespace := newTestResource("namespace", "test", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err := cs.CreateResource("profile", "testProfile", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "profile/testProfile created") {
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
	_, err = ds.DeleteResource("profile", "", "test")
	if err == nil || !strings.Contains(err.Error(), "no name was specified") {
		t.Errorf("expected missing name error, got %v", err)
	}

	// Test missing namespace
	_, err = ds.DeleteResource("profile", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	// Test successful deletion
	namespace := newTestResource("namespace", "test", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	profile := newTestResource("profile", "testProfile", "test")
	_, err = repo.Create(profile)
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
	_, err := ds.DescribeResource("profile", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	// Test successful description
	namespace := newTestResource("namespace", "testNS", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	profile := newTestResource("profile", "testProfile", "testNS")
	_, err = repo.Create(profile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err2 := ds.DescribeResource("profile", "testProfile", "testNS")
	if err2 != nil {
		t.Fatalf("unexpected error: %v", err2)
	}
	if !strings.Contains(msg, "testProfile") {
		t.Errorf("unexpected description output: %s", msg)
	}
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
	_, err = gs.GetResources("profile", "testName", "")
	if err == nil || !strings.Contains(err.Error(), "cannot be retrieved by name across all namespaces") {
		t.Errorf("expected missing namespace error, got %v", err)
	}

	// Test successful retrieval
	namespace := newTestResource("namespace", "testNS", "")
	_, err = repo.Create(namespace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	profile := newTestResource("profile", "testProfile", "testNS")
	_, err = repo.Create(profile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg, err := gs.GetResources("profile", "", "testNS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "testProfile") {
		t.Errorf("unexpected retrieval output: %s", msg)
	}
}

// Test for ResourceServiceImpl
func TestResourceServiceImpl(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	rs := service.NewResourceService(repo, cmd)

	// Test all delegations
	_, _ = rs.CreateResource("profile", "testProfile", "testNS")
	_, _ = rs.DeleteResource("profile", "testProfile", "testNS")
	_, _ = rs.GetResources("profile", "testProfile", "testNS")
	_, _ = rs.DescribeResource("profile", "testProfile", "testNS")
}
