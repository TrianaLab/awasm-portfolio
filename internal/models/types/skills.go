package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Skill struct {
	Competence  string
	Proficiency string
}

type Skills struct {
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	Skills            []Skill
	CreationTimestamp time.Time
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
