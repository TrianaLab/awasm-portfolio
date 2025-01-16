package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Skill struct {
	Category    string `yaml:"Category,omitempty"`
	Competence  string `yaml:"Competence,omitempty"`
	Proficiency string `yaml:"Proficiency,omitempty"`
}

type Skills struct {
	Kind              string                `yaml:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty"`
	Skills            []Skill               `yaml:"Skills,omitempty"`
}

func (s *Skills) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*s).Name()) }
func (s *Skills) GetName() string                               { return s.Name }
func (s *Skills) SetName(name string)                           { s.Name = name }
func (s *Skills) GetNamespace() string                          { return s.Namespace }
func (s *Skills) SetNamespace(namespace string)                 { s.Namespace = namespace }
func (s *Skills) GetOwnerReference() models.OwnerReference      { return s.OwnerRef }
func (s *Skills) SetOwnerReference(owner models.OwnerReference) { s.OwnerRef = owner }
func (s *Skills) GetID() string {
	return strings.ToLower(s.GetKind() + ":" + s.Name + ":" + s.Namespace)
}
func (s *Skills) GetCreationTimestamp() time.Time          { return s.CreationTimestamp }
func (s *Skills) SetCreationTimestamp(timestamp time.Time) { s.CreationTimestamp = timestamp }
