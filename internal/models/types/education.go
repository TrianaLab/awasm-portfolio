package types

import "awasm-portfolio/internal/models"

type Course struct {
	Title       string
	Institution string
	Duration    string
}

type Education struct {
	Name      string
	Namespace string
	OwnerRefs []models.OwnerReference
	Courses   []Course
}

func (e *Education) GetName() string                                   { return e.Name }
func (e *Education) SetName(name string)                               { e.Name = name }
func (e *Education) GetNamespace() string                              { return e.Namespace }
func (e *Education) SetNamespace(namespace string)                     { e.Namespace = namespace }
func (e *Education) GetOwnerReferences() []models.OwnerReference       { return e.OwnerRefs }
func (e *Education) SetOwnerReferences(owners []models.OwnerReference) { e.OwnerRefs = owners }
