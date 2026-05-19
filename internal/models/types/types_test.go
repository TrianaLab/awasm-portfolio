package types_test

import (
	"testing"
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
)

// resourceCase exercises the full Resource interface surface for every
// concrete type in the types package. The Namespace type is excluded from
// the SetNamespace / SetOwnerReference assertions because its setters are
// intentional no-ops (namespaces live at the cluster scope).
type resourceCase struct {
	name             string
	resource         models.Resource
	expectedKind     string
	settersAreNoOp   bool // true for Namespace, whose Set{Namespace,OwnerReference} ignore input
	getNamespaceIsNs bool // true for Namespace, whose GetNamespace() returns ""
}

func newCases() []resourceCase {
	return []resourceCase{
		{name: "Award", resource: &types.Award{}, expectedKind: "award"},
		{name: "Basics", resource: &types.Basics{}, expectedKind: "basics"},
		{name: "Certificate", resource: &types.Certificate{}, expectedKind: "certificate"},
		{name: "Education", resource: &types.Education{}, expectedKind: "education"},
		{name: "Interest", resource: &types.Interest{}, expectedKind: "interest"},
		{name: "Language", resource: &types.Language{}, expectedKind: "language"},
		{name: "Project", resource: &types.Project{}, expectedKind: "project"},
		{name: "Publication", resource: &types.Publication{}, expectedKind: "publication"},
		{name: "Reference", resource: &types.Reference{}, expectedKind: "reference"},
		{name: "Resume", resource: &types.Resume{}, expectedKind: "resume"},
		{name: "Skill", resource: &types.Skill{}, expectedKind: "skill"},
		{name: "Volunteer", resource: &types.Volunteer{}, expectedKind: "volunteer"},
		{name: "Work", resource: &types.Work{}, expectedKind: "work"},
		{
			name: "Namespace", resource: &types.Namespace{},
			expectedKind:   "namespace",
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

			if got := r.GetKind(); got != tc.expectedKind {
				t.Errorf("GetKind() = %q, want %q", got, tc.expectedKind)
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
