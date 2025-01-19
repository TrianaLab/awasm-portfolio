package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"testing"
	"time"
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

	// Test setting fields
	resource.SetName("UpdatedName")
	assertEqual(t, "SetName", "UpdatedName", resource.GetName())
	resource.SetNamespace("UpdatedNamespace")
	assertEqual(t, "SetNamespace", "UpdatedNamespace", resource.GetNamespace())

	// Test creation timestamp
	now := time.Now()
	resource.SetCreationTimestamp(now)
	assertEqual(t, "SetCreationTimestamp", now, resource.GetCreationTimestamp())
}

// TestCertifications tests the Certifications struct.
func TestCertifications(t *testing.T) {
	resource := &Certifications{
		Kind:      "certifications",
		Name:      "MyCertifications",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Certifications: []Certification{
			{Description: "Cert 1", Link: "http://cert1.com"},
		},
	}

	testResourceMethods(t, resource, "certifications", "MyCertifications", "default")
}

// TestContact tests the Contact struct.
func TestContact(t *testing.T) {
	resource := &Contact{
		Kind:      "contact",
		Name:      "JohnDoe",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Email:     "john@example.com",
		Linkedin:  "john-linkedin",
		Github:    "john-github",
	}

	testResourceMethods(t, resource, "contact", "JohnDoe", "default")
}

// TestContributions tests the Contributions struct.
func TestContributions(t *testing.T) {
	resource := &Contributions{
		Kind:      "contributions",
		Name:      "MyContributions",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Contributions: []Contribution{
			{Project: "Project 1", Description: "Desc 1", Link: "http://project1.com"},
		},
	}

	testResourceMethods(t, resource, "contributions", "MyContributions", "default")
}

// TestEducation tests the Education struct.
func TestEducation(t *testing.T) {
	resource := &Education{
		Kind:      "education",
		Name:      "MyEducation",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Courses: []Course{
			{Title: "Course 1", Institution: "Institution 1", Duration: "6 months"},
		},
	}

	testResourceMethods(t, resource, "education", "MyEducation", "default")
}

// TestExperience tests the Experience struct.
func TestExperience(t *testing.T) {
	resource := &Experience{
		Kind:      "experience",
		Name:      "MyExperience",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Jobs: []Job{
			{Title: "Job 1", Description: "Job Desc 1", Company: "Company 1", Duration: "1 year"},
		},
	}

	testResourceMethods(t, resource, "experience", "MyExperience", "default")
}

// TestNamespace tests the Namespace struct.
func TestNamespace(t *testing.T) {
	resource := &Namespace{
		Kind:     "namespace",
		Name:     "default",
		OwnerRef: newOwnerReference(),
	}

	assertEqual(t, "GetKind", "namespace", resource.GetKind())
	assertEqual(t, "GetName", "default", resource.GetName())
	assertEqual(t, "GetNamespace", "", resource.GetNamespace())
	assertEqual(t, "GetID", strings.ToLower("namespace:default:"), resource.GetID())

	// Test setting fields
	resource.SetName("UpdatedName")
	assertEqual(t, "SetName", "UpdatedName", resource.GetName())
	resource.SetNamespace("IgnoredNamespace")
	assertEqual(t, "SetNamespace", "", resource.GetNamespace())

	// Test creation timestamp
	now := time.Now()
	resource.SetCreationTimestamp(now)
	assertEqual(t, "SetCreationTimestamp", now, resource.GetCreationTimestamp())
}

// TestProfile tests the Profile struct.
func TestProfile(t *testing.T) {
	resource := &Profile{
		Kind:      "profile",
		Name:      "MyProfile",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Contributions: Contributions{
			Kind:      "contributions",
			Name:      "MyContributions",
			Namespace: "default",
		},
		Experience: Experience{
			Kind:      "experience",
			Name:      "MyExperience",
			Namespace: "default",
		},
		Certifications: Certifications{
			Kind:      "certifications",
			Name:      "MyCertifications",
			Namespace: "default",
		},
		Skills: Skills{
			Kind:      "skills",
			Name:      "MySkills",
			Namespace: "default",
		},
		Contact: Contact{
			Kind:      "contact",
			Name:      "JohnDoe",
			Namespace: "default",
		},
	}

	testResourceMethods(t, resource, "profile", "MyProfile", "default")
}

// TestSkills tests the Skills struct.
func TestSkills(t *testing.T) {
	resource := &Skills{
		Kind:      "skills",
		Name:      "MySkills",
		Namespace: "default",
		OwnerRef:  newOwnerReference(),
		Skills: []Skill{
			{Category: "Programming", Competence: "Expert", Proficiency: "5/5"},
		},
	}

	testResourceMethods(t, resource, "skills", "MySkills", "default")
}
