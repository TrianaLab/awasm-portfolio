package types_test

import (
	"testing"
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
)

// resourceCase exercises the full Resource interface surface for every
// concrete type. Namespace overrides Set{Namespace,OwnerReference} so its
// flags differ.
type resourceCase struct {
	name             string
	resource         models.Resource
	kind             string
	settersAreNoOp   bool // Namespace ignores SetNamespace + SetOwnerReference
	getNamespaceIsNs bool // Namespace.GetNamespace() returns ""
}

func newCases() []resourceCase {
	with := func(kind string) models.Meta { return models.Meta{Kind: kind} }
	return []resourceCase{
		{name: "Award", resource: &types.Award{Meta: with("award")}, kind: "award"},
		{name: "Basics", resource: &types.Basics{Meta: with("basics")}, kind: "basics"},
		{name: "Certificate", resource: &types.Certificate{Meta: with("certificate")}, kind: "certificate"},
		{name: "Education", resource: &types.Education{Meta: with("education")}, kind: "education"},
		{name: "Interest", resource: &types.Interest{Meta: with("interest")}, kind: "interest"},
		{name: "Language", resource: &types.Language{Meta: with("language")}, kind: "language"},
		{name: "Project", resource: &types.Project{Meta: with("project")}, kind: "project"},
		{name: "Publication", resource: &types.Publication{Meta: with("publication")}, kind: "publication"},
		{name: "Reference", resource: &types.Reference{Meta: with("reference")}, kind: "reference"},
		{name: "Resume", resource: &types.Resume{Meta: with("resume")}, kind: "resume"},
		{name: "Skill", resource: &types.Skill{Meta: with("skill")}, kind: "skill"},
		{name: "Volunteer", resource: &types.Volunteer{Meta: with("volunteer")}, kind: "volunteer"},
		{name: "Work", resource: &types.Work{Meta: with("work")}, kind: "work"},
		{
			name: "Namespace", resource: &types.Namespace{Meta: with("namespace")},
			kind:           "namespace",
			settersAreNoOp: true, getNamespaceIsNs: true,
		},
	}
}

func TestResourceInterface(t *testing.T) {
	owner := models.OwnerReference{Kind: "resume", Name: "owner", Namespace: "default"}
	ts := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)

	for _, tc := range newCases() {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.resource

			if got := r.GetKind(); got != tc.kind {
				t.Errorf("GetKind() = %q, want %q", got, tc.kind)
			}

			r.SetName("widget")
			if got := r.GetName(); got != "widget" {
				t.Errorf("after SetName(%q): GetName() = %q", "widget", got)
			}

			r.SetNamespace("ns-a")
			if tc.getNamespaceIsNs {
				if got := r.GetNamespace(); got != "" {
					t.Errorf("Namespace.GetNamespace() = %q, want \"\"", got)
				}
			} else {
				if got := r.GetNamespace(); got != "ns-a" {
					t.Errorf("after SetNamespace(%q): GetNamespace() = %q", "ns-a", got)
				}
			}

			r.SetOwnerReference(owner)
			gotOwner := r.GetOwnerReference()
			if tc.settersAreNoOp {
				if gotOwner.Name != "" {
					t.Errorf("Namespace.SetOwnerReference should be a no-op, got %+v", gotOwner)
				}
			} else {
				if gotOwner.Name != owner.Name {
					t.Errorf("GetOwnerReference().Name = %q, want %q", gotOwner.Name, owner.Name)
				}
			}

			id := r.GetID()
			if id == "" {
				t.Error("GetID() returned empty string")
			}

			r.SetCreationTimestamp(ts)
			if got := r.GetCreationTimestamp(); !got.Equal(ts) {
				t.Errorf("after SetCreationTimestamp: GetCreationTimestamp() = %v, want %v", got, ts)
			}
		})
	}
}
