package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
)

type Contribution struct {
	Project     string
	Description string
	Link        string
}

type Contributions struct {
	Name          string
	Namespace     string
	OwnerRef      models.OwnerReference
	Contributions []Contribution
}

func (c *Contributions) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*c).Name()) }
func (c *Contributions) GetName() string                               { return c.Name }
func (c *Contributions) SetName(name string)                           { c.Name = name }
func (c *Contributions) GetNamespace() string                          { return c.Namespace }
func (c *Contributions) SetNamespace(namespace string)                 { c.Namespace = namespace }
func (c *Contributions) GetOwnerReference() models.OwnerReference      { return c.OwnerRef }
func (c *Contributions) SetOwnerReference(owner models.OwnerReference) { c.OwnerRef = owner }
func (c *Contributions) GetID() string {
	return strings.ToLower(c.GetKind() + ":" + c.Name + ":" + c.Namespace)
}
