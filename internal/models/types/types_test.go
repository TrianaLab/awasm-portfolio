package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"testing"
)

// Helper to create an OwnerReference for testing.
func newOwnerReference() models.OwnerReference {
	return models.OwnerReference{
		Kind:      "namespace",
		Name:      "default",
		Namespace: "",
	}
}

// Helper to assert values.
func assertEqual(t *testing.T, field string, expected, actual any) {
	if expected != actual {
		t.Errorf("expected %s to be %v, got %v", field, expected, actual)
	}
}

// Test common methods for all types implementing Resource interface.
func testResourceMethods(t *testing.T, resource models.Resource, kind, name, namespace string) {
	assertEqual(t, "GetKind", kind, resource.GetKind())
	assertEqual(t, "GetName", name, resource.GetName())
	assertEqual(t, "GetNamespace", namespace, resource.GetNamespace())
	assertEqual(t, "GetID", strings.ToLower(kind+":"+name+":"+namespace), resource.GetID())
}

// Test each type implementation
func TestCertifications(t *testing.T) {
	resource := &Certifications{
		Kind:      "certifications",
		Name:      "MyCertifications",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "certifications", "MyCertifications", "default")

	// Test setting fields
	resource.SetName("UpdatedCertifications")
	assertEqual(t, "SetName", "UpdatedCertifications", resource.GetName())
	resource.SetNamespace("updated-namespace")
	assertEqual(t, "SetNamespace", "updated-namespace", resource.GetNamespace())
}

func TestContact(t *testing.T) {
	resource := &Contact{
		Kind:      "contact",
		Name:      "JohnDoe",
		Namespace: "default",
		Email:     "john@example.com",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "contact", "JohnDoe", "default")

	// Test setting fields
	resource.SetName("JaneDoe")
	assertEqual(t, "SetName", "JaneDoe", resource.GetName())
	resource.SetNamespace("updated-namespace")
	assertEqual(t, "SetNamespace", "updated-namespace", resource.GetNamespace())
}

func TestContributions(t *testing.T) {
	resource := &Contributions{
		Kind:      "contributions",
		Name:      "MyContributions",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "contributions", "MyContributions", "default")
}

func TestEducation(t *testing.T) {
	resource := &Education{
		Kind:      "education",
		Name:      "MyEducation",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "education", "MyEducation", "default")
}

func TestExperience(t *testing.T) {
	resource := &Experience{
		Kind:      "experience",
		Name:      "MyExperience",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "experience", "MyExperience", "default")
}

func TestNamespace(t *testing.T) {
	resource := &Namespace{
		Kind:     "namespace",
		Name:     "default",
		OwnerRef: newOwnerReference(),
	}

	testResourceMethods(t, resource, "namespace", "default", "")

	// Test setting fields
	resource.SetName("UpdatedNamespace")
	assertEqual(t, "SetName", "UpdatedNamespace", resource.GetName())
}

func TestProfile(t *testing.T) {
	resource := &Profile{
		Kind:      "profile",
		Name:      "MyProfile",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "profile", "MyProfile", "default")
}

func TestSkills(t *testing.T) {
	resource := &Skills{
		Kind:      "skills",
		Name:      "MySkills",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
	}

	testResourceMethods(t, resource, "skills", "MySkills", "default")
}
