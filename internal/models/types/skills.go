package types

import "awasm-portfolio/internal/models"

type Skill struct {
	Name        string
	Proficiency string
}

type Skills struct {
	Name      string
	Namespace string
	OwnerRefs []models.OwnerReference
	Skills    []Skill
}

func (s *Skills) GetName() string                                   { return s.Name }
func (s *Skills) SetName(name string)                               { s.Name = name }
func (s *Skills) GetNamespace() string                              { return s.Namespace }
func (s *Skills) SetNamespace(namespace string)                     { s.Namespace = namespace }
func (s *Skills) GetOwnerReferences() []models.OwnerReference       { return s.OwnerRefs }
func (s *Skills) SetOwnerReferences(owners []models.OwnerReference) { s.OwnerRefs = owners }
