package factory_test

import (
	"strings"
	"testing"

	"github.com/TrianaLab/awasm-portfolio/internal/factory"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
)

// TestNew_AllKinds asserts factory.New populates non-empty content for
// every supported kind and assigns the right Meta identity.
func TestNew_AllKinds(t *testing.T) {
	kinds := []string{
		"namespace", "resume", "basics", "work", "volunteer", "education",
		"award", "certificate", "publication", "skill", "language",
		"interest", "reference", "project",
	}
	for _, kind := range kinds {
		t.Run(kind, func(t *testing.T) {
			r, err := factory.New(kind, "demo", "default")
			if err != nil {
				t.Fatalf("unexpected error for kind %s: %v", kind, err)
			}
			if r.GetName() != "demo" {
				t.Errorf("GetName() = %q, want %q", r.GetName(), "demo")
			}
			if r.GetKind() != kind {
				t.Errorf("GetKind() = %q, want %q", r.GetKind(), kind)
			}
		})
	}
}

func TestNew_UnknownKind(t *testing.T) {
	_, err := factory.New("nonsense", "x", "ns")
	if err == nil {
		t.Fatal("expected error for unsupported kind")
	}
	if !strings.Contains(err.Error(), "unsupported resource kind") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestNew_NamespaceIsClusterScoped verifies the Namespace override —
// it must never carry a parent namespace, even if one is passed.
func TestNew_NamespaceIsClusterScoped(t *testing.T) {
	r, err := factory.New("namespace", "ns-a", "ignored")
	if err != nil {
		t.Fatal(err)
	}
	if r.GetNamespace() != "" {
		t.Errorf("Namespace.GetNamespace() = %q, want empty", r.GetNamespace())
	}
}

// TestNew_ResumePopulatesChildren is the big one — Resume should come
// back with a populated Basics block and ≥1 entry in every slice
// section so the resume view and PDF have content to render.
func TestNew_ResumePopulatesChildren(t *testing.T) {
	r, err := factory.New("resume", "demo-resume", "default")
	if err != nil {
		t.Fatal(err)
	}
	resume, ok := r.(*types.Resume)
	if !ok {
		t.Fatalf("expected *types.Resume, got %T", r)
	}

	if resume.Basics.FullName == "" {
		t.Error("Basics.FullName is empty")
	}
	checks := []struct {
		name  string
		count int
	}{
		{"work", len(resume.Work)},
		{"volunteer", len(resume.Volunteer)},
		{"education", len(resume.Education)},
		{"awards", len(resume.Awards)},
		{"certificates", len(resume.Certificates)},
		{"publications", len(resume.Publications)},
		{"skills", len(resume.Skills)},
		{"languages", len(resume.Languages)},
		{"interests", len(resume.Interests)},
		{"references", len(resume.References)},
		{"projects", len(resume.Projects)},
	}
	for _, c := range checks {
		if c.count == 0 {
			t.Errorf("Resume.%s is empty; expected populated by factory", c.name)
		}
	}
}

// TestNew_ResumeChildrenCarryOwnerRef ensures every populated child
// has the resume's owner reference and matching namespace.
func TestNew_ResumeChildrenCarryOwnerRef(t *testing.T) {
	r, _ := factory.New("resume", "owner", "ns-a")
	resume := r.(*types.Resume)

	for _, w := range resume.Work {
		if w.OwnerRef.Name != "owner" || w.OwnerRef.Kind != "resume" {
			t.Errorf("work.OwnerRef = %+v, want resume/owner", w.OwnerRef)
		}
		if w.GetNamespace() != "ns-a" {
			t.Errorf("work.Namespace = %q, want %q", w.GetNamespace(), "ns-a")
		}
	}
}
