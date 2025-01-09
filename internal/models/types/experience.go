package types

import "awasm-portfolio/internal/models"

type Job struct {
	Title       string
	Description string
	Company     string
	Duration    string
}

type Experience struct {
	Name      string
	Namespace string
	OwnerRefs []models.OwnerReference
	Jobs      []Job
}

func (e *Experience) GetName() string                                   { return e.Name }
func (e *Experience) SetName(name string)                               { e.Name = name }
func (e *Experience) GetNamespace() string                              { return e.Namespace }
func (e *Experience) SetNamespace(namespace string)                     { e.Namespace = namespace }
func (e *Experience) GetOwnerReferences() []models.OwnerReference       { return e.OwnerRefs }
func (e *Experience) SetOwnerReferences(owners []models.OwnerReference) { e.OwnerRefs = owners }
