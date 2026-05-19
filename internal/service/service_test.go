package service_test

import (
	"strings"
	"testing"
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
	"github.com/TrianaLab/awasm-portfolio/internal/service"
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

func newTestResource(kind, name, namespace string) models.Resource {
	return &dummyResource{Kind: kind, Name: name, Namespace: namespace}
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

// ---------------- Create ----------------

func TestCreate_InvalidKind(t *testing.T) {
	_, err := service.Create(repository.NewInMemoryRepository(), "invalidKind", "x", "ns")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestCreate_MissingNamespace(t *testing.T) {
	_, err := service.Create(repository.NewInMemoryRepository(), "resume", "x", "nonexistentNS")
	assertErrContains(t, err, "namespace 'nonexistentNS' not found")
}

func TestCreate_Success(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "test", ""))

	msg, err := service.Create(repo, "resume", "testResume", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "resume/testResume created") {
		t.Errorf("unexpected success message: %s", msg)
	}
}

func TestCreate_Duplicate(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "dup", ""))

	_, err := service.Create(repo, "namespace", "dup", "")
	if err == nil {
		t.Error("expected error creating duplicate namespace")
	}
}

func TestCreate_AllKinds(t *testing.T) {
	// Exercise every branch of newResource so it stays 100% covered.
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "all-kinds", ""))

	for _, kind := range []string{
		"resume", "basics", "work", "volunteer", "education",
		"award", "certificate", "publication", "skill", "language",
		"interest", "reference", "project",
	} {
		if _, err := service.Create(repo, kind, kind+"-instance", "all-kinds"); err != nil {
			t.Errorf("Create(%s): %v", kind, err)
		}
	}
}

// "all" passes NormalizeKind (returns empty kind) but is rejected by
// newResource — exercises the post-NormalizeKind error path in Create.
func TestCreate_AllRejected(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "ns", ""))
	_, err := service.Create(repo, "all", "x", "ns")
	assertErrContains(t, err, "unsupported resource kind")
}

// ---------------- Delete ----------------

func TestDelete_InvalidKind(t *testing.T) {
	_, err := service.Delete(repository.NewInMemoryRepository(), "invalidKind", "x", "test")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestDelete_AllKindRejected(t *testing.T) {
	_, err := service.Delete(repository.NewInMemoryRepository(), "all", "", "default")
	assertErrContains(t, err, "you must specify only one resource")
}

func TestDelete_MissingName(t *testing.T) {
	_, err := service.Delete(repository.NewInMemoryRepository(), "resume", "", "test")
	assertErrContains(t, err, "no name was specified")
}

func TestDelete_MissingNamespace(t *testing.T) {
	_, err := service.Delete(repository.NewInMemoryRepository(), "resume", "x", "")
	assertErrContains(t, err, "cannot be retrieved by name across all namespaces")
}

func TestDelete_Namespace(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "test", ""))
	mustCreate(t, repo, newTestResource("resume", "testResume", "test"))

	msg, err := service.Delete(repo, "namespace", "test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "namespace/test in namespace '' deleted") {
		t.Errorf("unexpected success message: %s", msg)
	}
}

func TestDelete_NamespaceNotFound(t *testing.T) {
	_, err := service.Delete(repository.NewInMemoryRepository(), "namespace", "inexistent", "")
	assertErrContains(t, err, "namespace/inexistent not found in namespace ''")
}

func TestDelete_WithOwnerReferences(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "test-ns", ""))
	mustCreate(t, repo, newTestResource("resume", "parent", "test-ns"))
	parentOwner := models.OwnerReference{Kind: "resume", Name: "parent", Namespace: "test-ns"}
	mustCreate(t, repo, &dummyResource{Kind: "work", Name: "child-work", Namespace: "test-ns", ownerRef: parentOwner})
	mustCreate(t, repo, &dummyResource{Kind: "education", Name: "child-edu", Namespace: "test-ns", ownerRef: parentOwner})

	msg, err := service.Delete(repo, "resume", "parent", "test-ns")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"resume/parent", "work/child-work", "education/child-edu"} {
		if !strings.Contains(msg, want) {
			t.Errorf("delete message missing %q: %s", want, msg)
		}
	}
}

func TestDelete_NonexistentResume(t *testing.T) {
	_, err := service.Delete(repository.NewInMemoryRepository(), "resume", "nonexistent", "nonexistent")
	assertErrContains(t, err, "not found")
}

// ---------------- Get ----------------

func TestGet_MissingNamespace(t *testing.T) {
	_, err := service.Get(repository.NewInMemoryRepository(), "resume", "testName", "", "json")
	assertErrContains(t, err, "cannot be retrieved by name across all namespaces")
}

func TestGet_Success(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "testNS", ""))
	mustCreate(t, repo, newTestResource("resume", "testResume", "testNS"))

	msg, err := service.Get(repo, "resume", "", "testNS", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "testResume") {
		t.Errorf("unexpected output: %s", msg)
	}
}

func TestGet_InvalidKind(t *testing.T) {
	_, err := service.Get(repository.NewInMemoryRepository(), "nonsense", "", "", "")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestGet_AllExcludesNamespaces(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	for _, r := range []struct{ kind, name, ns string }{
		{"namespace", "ns1", ""}, {"namespace", "ns2", ""},
		{"work", "work-1", "ns1"}, {"education", "edu-1", "ns1"},
	} {
		mustCreate(t, repo, newTestResource(r.kind, r.name, r.ns))
	}

	msg, err := service.Get(repo, "all", "", "", "json")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if strings.Contains(msg, "\"Kind\": \"namespace\"") {
		t.Error("namespaces should not be included in output")
	}
	if !strings.Contains(msg, "\"Kind\": \"work\"") {
		t.Error("work resources should be included in output")
	}
}

// ---------------- Describe ----------------

func TestDescribe_MissingNamespace(t *testing.T) {
	_, err := service.Describe(repository.NewInMemoryRepository(), "resume", "testName", "")
	assertErrContains(t, err, "cannot be retrieved by name across all namespaces")
}

func TestDescribe_Success(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "testNS", ""))
	mustCreate(t, repo, newTestResource("resume", "testResume", "testNS"))

	msg, err := service.Describe(repo, "resume", "testResume", "testNS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "testResume") {
		t.Errorf("unexpected description output: %s", msg)
	}
}

func TestDescribe_InvalidKind(t *testing.T) {
	_, err := service.Describe(repository.NewInMemoryRepository(), "nonsense", "", "")
	assertErrContains(t, err, "unsupported resource kind")
}

func TestDescribe_AllInNamespace(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "testNS", ""))
	for _, k := range []struct{ kind, name string }{
		{"work", "work-1"},
		{"education", "edu-1"},
		{"skill", "skill-1"},
	} {
		mustCreate(t, repo, newTestResource(k.kind, k.name, "testNS"))
	}

	for _, k := range []struct{ kind, name string }{
		{"work", "work-1"}, {"education", "edu-1"}, {"skill", "skill-1"},
	} {
		msg, err := service.Describe(repo, k.kind, "", "testNS")
		if err != nil {
			t.Errorf("describe %s: %v", k.kind, err)
		}
		if !strings.Contains(msg, k.name) {
			t.Errorf("expected %s in output, got: %s", k.name, msg)
		}
	}
}

func TestDescribe_AllWithNamespaceFilter(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	mustCreate(t, repo, newTestResource("namespace", "ns1", ""))
	mustCreate(t, repo, newTestResource("namespace", "ns2", ""))
	for _, r := range []struct{ kind, name, ns string }{
		{"work", "work-1", "ns1"}, {"work", "work-2", "ns2"},
		{"education", "edu-1", "ns1"}, {"education", "edu-2", "ns2"},
	} {
		mustCreate(t, repo, newTestResource(r.kind, r.name, r.ns))
	}

	msg, err := service.Describe(repo, "all", "", "ns1")
	if err != nil {
		t.Fatalf("error: %v", err)
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
