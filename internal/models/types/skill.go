package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Skill struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"name" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Level             string                `json:"level" yaml:"Level,omitempty"`
	Keywords          []string              `json:"keywords" yaml:"Keywords,omitempty"`
}

func (s *Skill) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*s).Name()) }
func (s *Skill) GetName() string                                { return s.Name }
func (s *Skill) SetName(name string)                            { s.Name = name }
func (s *Skill) GetNamespace() string                           { return s.Namespace }
func (s *Skill) SetNamespace(namespace string)                  { s.Namespace = namespace }
func (s *Skill) GetOwnerReference() models.OwnerReference       { return s.OwnerRef }
func (s *Skill) SetOwnerReference(owners models.OwnerReference) { s.OwnerRef = owners }
func (s *Skill) GetID() string {
	return strings.ToLower(s.GetKind() + ":" + s.Name + ":" + s.Namespace)
}
func (s *Skill) GetCreationTimestamp() time.Time          { return s.CreationTimestamp }
func (s *Skill) SetCreationTimestamp(timestamp time.Time) { s.CreationTimestamp = timestamp }
