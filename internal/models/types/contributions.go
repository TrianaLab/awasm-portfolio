package types

import "awasm-portfolio/internal/models"

type Contribution struct {
	Project     string
	Description string
	Link        string
}

type Contributions struct {
	Name          string
	Namespace     string
	OwnerRefs     []models.OwnerReference
	Contributions []Contribution
}

func (c *Contributions) GetName() string                                   { return c.Name }
func (c *Contributions) SetName(name string)                               { c.Name = name }
func (c *Contributions) GetNamespace() string                              { return c.Namespace }
func (c *Contributions) SetNamespace(namespace string)                     { c.Namespace = namespace }
func (c *Contributions) GetOwnerReferences() []models.OwnerReference       { return c.OwnerRefs }
func (c *Contributions) SetOwnerReferences(owners []models.OwnerReference) { c.OwnerRefs = owners }
