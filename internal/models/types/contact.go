package types

import "awasm-portfolio/internal/models"

type Contact struct {
	Name      string
	Namespace string
	OwnerRefs []models.OwnerReference
	Email     string
	LinkedIn  string
	GitHub    string
}

func (c *Contact) GetName() string                                   { return c.Name }
func (c *Contact) SetName(name string)                               { c.Name = name }
func (c *Contact) GetNamespace() string                              { return c.Namespace }
func (c *Contact) SetNamespace(namespace string)                     { c.Namespace = namespace }
func (c *Contact) GetOwnerReferences() []models.OwnerReference       { return c.OwnerRefs }
func (c *Contact) SetOwnerReferences(owners []models.OwnerReference) { c.OwnerRefs = owners }
