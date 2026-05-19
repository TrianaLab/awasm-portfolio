package service_test

import (
	"strings"
	"testing"
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/service"

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

func mustCreate(t *testing.T, repo *repository.InMemoryRepository, r models.Resource) {
	t.Helper()
	if _, err := repo.Create(r); err != nil {
		t.Fatalf("setup: create %s/%s: %v", r.GetKind(), r.GetName(), err)
	}
}

func assertErrContains(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", substr)
		return
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("expected error containing %q, got %v", substr, err)
	}
}

// ---------------- CreateService ----------------

func newCreateServiceFixture() (*repository.InMemoryRepository, *service.CreateService) {
	repo := repository.NewInMemoryRepository()
	return repo, service.NewCreateService(repo, newTestCommand())
}

func TestCreateService_InvalidKind(t *testing.T) {
	_, cs := newCreateServiceFixture()
	_, err := cs.CreateResource("invalidKind", "testName", "testNS")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestCreateService_MissingNamespace(t *testing.T) {
	_, cs := newCreateServiceFixture()
	_, err := cs.CreateResource("resume", "testName", "nonexistentNS")
	assertErrContains(t, err, "namespace 'nonexistentNS' not found")
}

func TestCreateService_Success(t *testing.T) {
	repo, cs := newCreateServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "test", ""))

	msg, err := cs.CreateResource("resume", "testResume", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "resume/testResume created") {
		t.Errorf("unexpected success message: %s", msg)
	}
}

func TestCreateService_DuplicateErrors(t *testing.T) {
	repo, cs := newCreateServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "dup", ""))

	_, err := cs.CreateResource("namespace", "dup", "")
	if err == nil {
		t.Error("expected error creating duplicate namespace")
	}
}

func TestCreateService_ResumeBasicsConflict(t *testing.T) {
	repo, cs := newCreateServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "basics-conflict-ns", ""))
	// Pre-create the basics resource that the resume would derive.
	mustCreate(t, repo, newTestResource("basics", "basics-conflict-basics", "basics-conflict-ns"))

	_, err := cs.CreateResource("resume", "basics-conflict", "basics-conflict-ns")
	assertErrContains(t, err, "failed to save basics")
}

func TestCreateService_ResumeNestedConflict(t *testing.T) {
	repo, cs := newCreateServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "nested-conflict-ns", ""))
	// First nested work entry the factory produces is "<name>-Work-0".
	mustCreate(t, repo, newTestResource("work", "nested-conflict-work-0", "nested-conflict-ns"))

	_, err := cs.CreateResource("resume", "nested-conflict", "nested-conflict-ns")
	assertErrContains(t, err, "failed to save")
}

func TestCreateService_ResumeWithNestedResources(t *testing.T) {
	repo, cs := newCreateServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "nested-test", ""))

	if _, err := cs.CreateResource("resume", "test-nested", "nested-test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	kinds := []string{
		"work", "education", "volunteer", "award", "skill",
		"language", "interest", "reference", "project", "basics",
	}
	for _, typ := range kinds {
		assertNestedChildren(t, repo, typ, "test-nested", "nested-test")
	}
}

func assertNestedChildren(t *testing.T, repo *repository.InMemoryRepository, kind, parent, namespace string) {
	t.Helper()
	resources, err := repo.List(kind, "", namespace)
	if err != nil {
		t.Errorf("error listing %s: %v", kind, err)
	}
	if len(resources) == 0 {
		t.Errorf("no %s resources created", kind)
	}
	for _, res := range resources {
		owner := res.GetOwnerReference()
		if owner.Kind != "resume" || owner.Name != parent || owner.Namespace != namespace {
			t.Errorf("%s/%s: wrong owner %+v", kind, res.GetName(), owner)
		}
		if res.GetNamespace() != namespace {
			t.Errorf("%s/%s: wrong namespace %q", kind, res.GetName(), res.GetNamespace())
		}
	}
}

// ---------------- DeleteService ----------------

func newDeleteServiceFixture() (*repository.InMemoryRepository, *service.DeleteService) {
	repo := repository.NewInMemoryRepository()
	return repo, service.NewDeleteService(repo, newTestCommand())
}

func TestDeleteService_InvalidKind(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("invalidKind", "testName", "test")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestDeleteService_MissingName(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("resume", "", "test")
	assertErrContains(t, err, "no name was specified")
}

func TestDeleteService_MissingNamespace(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("resume", "testName", "")
	assertErrContains(t, err, "cannot be retrieved by name across all namespaces")
}

func TestDeleteService_Namespace(t *testing.T) {
	repo, ds := newDeleteServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "test", ""))
	mustCreate(t, repo, newTestResource("resume", "testResume", "test"))

	msg, err := ds.DeleteResource("namespace", "test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "namespace/test in namespace '' deleted") {
		t.Errorf("unexpected success message: %s", msg)
	}
}

func TestDeleteService_NamespaceNotFound(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("namespace", "inexistent", "")
	assertErrContains(t, err, "namespace/inexistent not found in namespace ''")
}

func TestDeleteService_EmptyKind(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("", "test", "default")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestDeleteService_AllNoName(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("all", "", "default")
	assertErrContains(t, err, "you must specify only one resource")
}

func TestDeleteService_WithOwnerReferences(t *testing.T) {
	repo, ds := newDeleteServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "test-ns", ""))
	mustCreate(t, repo, newTestResource("resume", "parent-resume", "test-ns"))
	parentOwner := models.OwnerReference{Kind: "resume", Name: "parent-resume", Namespace: "test-ns"}
	mustCreate(t, repo, &dummyResource{Kind: "work", Name: "child-work", Namespace: "test-ns", ownerRef: parentOwner})
	mustCreate(t, repo, &dummyResource{Kind: "education", Name: "child-edu", Namespace: "test-ns", ownerRef: parentOwner})

	msg, err := ds.DeleteResource("resume", "parent-resume", "test-ns")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"resume/parent-resume", "work/child-work", "education/child-edu"} {
		if !strings.Contains(msg, want) {
			t.Errorf("delete message missing %q: %s", want, msg)
		}
	}
	assertEmptyList(t, repo, "work", "test-ns")
	assertEmptyList(t, repo, "education", "test-ns")
}

func assertEmptyList(t *testing.T, repo *repository.InMemoryRepository, kind, namespace string) {
	t.Helper()
	resources, _ := repo.List(kind, "", namespace)
	if len(resources) > 0 {
		t.Errorf("expected no %s remaining in %s, got %d", kind, namespace, len(resources))
	}
}

func TestDeleteService_NonexistentResume(t *testing.T) {
	_, ds := newDeleteServiceFixture()
	_, err := ds.DeleteResource("resume", "nonexistent", "nonexistent")
	assertErrContains(t, err, "not found")
}

func TestDeleteService_ChildThenNamespace(t *testing.T) {
	repo, ds := newDeleteServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "test-ns2", ""))
	mustCreate(t, repo, newTestResource("resume", "parent2", "test-ns2"))
	mustCreate(t, repo, &dummyResource{
		Kind: "work", Name: "child2", Namespace: "test-ns2",
		ownerRef: models.OwnerReference{Kind: "resume", Name: "parent2", Namespace: "test-ns2"},
	})

	msg, err := ds.DeleteResource("work", "child2", "test-ns2")
	if err != nil {
		t.Errorf("unexpected error deleting child: %v", err)
	}
	if !strings.Contains(msg, "work/child2") {
		t.Error("child deletion message not found")
	}

	msg, err = ds.DeleteResource("namespace", "test-ns2", "")
	if err != nil {
		t.Errorf("unexpected error deleting namespace: %v", err)
	}
	if !strings.Contains(msg, "namespace/test-ns2") {
		t.Error("namespace deletion message not found")
	}
}

// ---------------- DescribeService ----------------

func newDescribeServiceFixture() (*repository.InMemoryRepository, *service.DescribeService) {
	repo := repository.NewInMemoryRepository()
	return repo, service.NewDescribeService(repo, newTestCommand())
}

func TestDescribeService_MissingNamespace(t *testing.T) {
	_, ds := newDescribeServiceFixture()
	_, err := ds.DescribeResource("resume", "testName", "")
	assertErrContains(t, err, "cannot be retrieved by name across all namespaces")
}

func TestDescribeService_Success(t *testing.T) {
	repo, ds := newDescribeServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "testNS", ""))
	mustCreate(t, repo, newTestResource("resume", "testResume", "testNS"))

	msg, err := ds.DescribeResource("resume", "testResume", "testNS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "testResume") {
		t.Errorf("unexpected description output: %s", msg)
	}
}

func TestDescribeService_AllInNamespace(t *testing.T) {
	repo, ds := newDescribeServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "testNS", ""))
	kinds := []struct{ kind, name string }{
		{"work", "work-1"},
		{"education", "edu-1"},
		{"skill", "skill-1"},
	}
	for _, k := range kinds {
		mustCreate(t, repo, newTestResource(k.kind, k.name, "testNS"))
	}

	for _, k := range kinds {
		msg, err := ds.DescribeResource(k.kind, "", "testNS")
		if err != nil {
			t.Errorf("error describing all %s: %v", k.kind, err)
		}
		if !strings.Contains(msg, k.name) {
			t.Errorf("expected %s in output, got: %s", k.name, msg)
		}
	}
}

func TestDescribeService_AllWithNamespaceFilter(t *testing.T) {
	repo, ds := newDescribeServiceFixture()
	mustCreate(t, repo, newTestResource("namespace", "ns1", ""))
	mustCreate(t, repo, newTestResource("namespace", "ns2", ""))
	for _, r := range []struct{ kind, name, ns string }{
		{"work", "work-1", "ns1"}, {"work", "work-2", "ns2"},
		{"education", "edu-1", "ns1"}, {"education", "edu-2", "ns2"},
	} {
		mustCreate(t, repo, newTestResource(r.kind, r.name, r.ns))
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
}

// ---------------- GetService ----------------

func newGetServiceFixture(t *testing.T) (*repository.InMemoryRepository, *cobra.Command, *service.GetService) {
	t.Helper()
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	if err := cmd.Flags().Set("output", "json"); err != nil {
		t.Fatalf("flag setup: %v", err)
	}
	return repo, cmd, service.NewGetService(repo, cmd)
}

func TestGetService_MissingNamespace(t *testing.T) {
	_, _, gs := newGetServiceFixture(t)
	_, err := gs.GetResources("resume", "testName", "")
	assertErrContains(t, err, "cannot be retrieved by name across all namespaces")
}

func TestGetService_Success(t *testing.T) {
	repo, _, gs := newGetServiceFixture(t)
	mustCreate(t, repo, newTestResource("namespace", "testNS", ""))
	mustCreate(t, repo, newTestResource("resume", "testResume", "testNS"))

	msg, err := gs.GetResources("resume", "", "testNS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "testResume") {
		t.Errorf("unexpected retrieval output: %s", msg)
	}
}

func TestGetService_AllAcrossNamespaces(t *testing.T) {
	repo, cmd, gs := newGetServiceFixture(t)
	for _, r := range []struct{ kind, name, ns string }{
		{"namespace", "ns1", ""}, {"namespace", "ns2", ""},
		{"work", "work-1", "ns1"}, {"education", "edu-1", "ns1"},
	} {
		mustCreate(t, repo, newTestResource(r.kind, r.name, r.ns))
	}
	if err := cmd.Flags().Set("output", "json"); err != nil {
		t.Fatalf("flag setup: %v", err)
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
}

// ---------------- ResourceService aggregate ----------------

func TestResourceServiceImpl(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	cmd := newTestCommand()
	rs := service.NewResourceService(repo, cmd)

	_, _ = rs.CreateResource("resume", "testResume", "testNS")
	_, _ = rs.DeleteResource("resume", "testResume", "testNS")
	_, _ = rs.GetResources("resume", "testResume", "testNS")
	_, _ = rs.DescribeResource("resume", "testResume", "testNS")
}
