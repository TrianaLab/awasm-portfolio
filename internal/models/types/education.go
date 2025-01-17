package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"time"
)

type Course struct {
	Title       string `yaml:"Title,omitempty"`
	Institution string `yaml:"Institution,omitempty"`
	Duration    string `yaml:"Duration,omitempty"`
}

type Education struct {
	Kind              string                `yaml:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty"`
	Courses           []Course              `yaml:"Courses,omitempty"`
}

func (e *Education) GetKind() string                               { return "education" }
func (e *Education) GetName() string                               { return e.Name }
func (e *Education) SetName(name string)                           { e.Name = name }
func (e *Education) GetNamespace() string                          { return e.Namespace }
func (e *Education) SetNamespace(namespace string)                 { e.Namespace = namespace }
func (e *Education) GetOwnerReference() models.OwnerReference      { return e.OwnerRef }
func (e *Education) SetOwnerReference(owner models.OwnerReference) { e.OwnerRef = owner }
func (e *Education) GetID() string {
	return strings.ToLower(e.GetKind() + ":" + e.Name + ":" + e.Namespace)
}
func (e *Education) GetCreationTimestamp() time.Time          { return e.CreationTimestamp }
func (e *Education) SetCreationTimestamp(timestamp time.Time) { e.CreationTimestamp = timestamp }
