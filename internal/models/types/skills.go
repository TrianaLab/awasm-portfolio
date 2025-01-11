package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
)

type Skill struct {
	Name        string
	Proficiency string
}

type Skills struct {
	Name      string
	Namespace string
	OwnerRef  models.OwnerReference
	Skills    []Skill
}

func (s *Skills) GetKind() string                               { return reflect.TypeOf(*s).Name() }
func (s *Skills) GetName() string                               { return s.Name }
func (s *Skills) SetName(name string)                           { s.Name = name }
func (s *Skills) GetNamespace() string                          { return s.Namespace }
func (s *Skills) SetNamespace(namespace string)                 { s.Namespace = namespace }
func (s *Skills) GetOwnerReference() models.OwnerReference      { return s.OwnerRef }
func (s *Skills) SetOwnerReference(owner models.OwnerReference) { s.OwnerRef = owner }
