package factory_test

import (
	"awasm-portfolio/internal/factory"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		kind         string
		data         map[string]interface{}
		expectNil    bool
		expectedKind string
	}{
		{
			name: "create namespace",
			kind: "namespace",
			data: map[string]interface{}{
				"name":      "test-namespace",
				"namespace": "",
			},
			expectNil:    false,
			expectedKind: "namespace",
		},
		{
			name: "create resume",
			kind: "resume",
			data: map[string]interface{}{
				"name":      "test-resume",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "resume",
		},
		{
			name: "create work",
			kind: "work",
			data: map[string]interface{}{
				"name":      "test-work",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "work",
		},
		{
			name: "create basics",
			kind: "basics",
			data: map[string]interface{}{
				"name":      "test-basics",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "basics",
		},
		{
			name: "create volunteer",
			kind: "volunteer",
			data: map[string]interface{}{
				"name":      "test-volunteer",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "volunteer",
		},
		{
			name: "create skill",
			kind: "skill",
			data: map[string]interface{}{
				"name":      "test-skill",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "skill",
		},
		{
			name: "create reference",
			kind: "reference",
			data: map[string]interface{}{
				"name":      "test-reference",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "reference",
		},
		{
			name: "create publication",
			kind: "publication",
			data: map[string]interface{}{
				"name":      "test-publication",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "publication",
		},
		{
			name: "create project",
			kind: "project",
			data: map[string]interface{}{
				"name":      "test-project",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "project",
		},
		{
			name: "create language",
			kind: "language",
			data: map[string]interface{}{
				"name":      "test-language",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "language",
		},
		{
			name: "create interest",
			kind: "interest",
			data: map[string]interface{}{
				"name":      "test-interest",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "interest",
		},
		{
			name: "create education",
			kind: "education",
			data: map[string]interface{}{
				"name":      "test-education",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "education",
		},
		{
			name: "create certificate",
			kind: "certificate",
			data: map[string]interface{}{
				"name":      "test-certificate",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "certificate",
		},
		{
			name: "create award",
			kind: "award",
			data: map[string]interface{}{
				"name":      "test-award",
				"namespace": "default",
			},
			expectNil:    false,
			expectedKind: "award",
		},
		{
			name: "invalid kind",
			kind: "invalid",
			data: map[string]interface{}{
				"name":      "test",
				"namespace": "default",
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := factory.NewResourceFactory()
			resource := f.Create(tt.kind, tt.data)

			if tt.expectNil {
				if resource != nil {
					t.Error("expected nil resource but got one")
				}
				return
			}

			if resource == nil {
				t.Fatal("expected resource but got nil")
			}

			if resource.GetKind() != tt.expectedKind {
				t.Errorf("expected kind %q, got %q", tt.expectedKind, resource.GetKind())
			}
		})
	}
}
