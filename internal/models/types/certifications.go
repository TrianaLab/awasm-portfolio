package types

import "awasm-portfolio/internal/models"

type Certification struct {
	Description string
	Link        string
}

type Certifications struct {
	Name           string
	Namespace      string
	OwnerRefs      []models.OwnerReference
	Certifications []Certification
}

func (c *Certifications) GetName() string                                   { return c.Name }
func (c *Certifications) SetName(name string)                               { c.Name = name }
func (c *Certifications) GetNamespace() string                              { return c.Namespace }
func (c *Certifications) SetNamespace(namespace string)                     { c.Namespace = namespace }
func (c *Certifications) GetOwnerReferences() []models.OwnerReference       { return c.OwnerRefs }
func (c *Certifications) SetOwnerReferences(owners []models.OwnerReference) { c.OwnerRefs = owners }
