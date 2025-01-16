package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Certification struct {
	Description string
	Link        string
}

type Certifications struct {
	Kind              string
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	Certifications    []Certification
	CreationTimestamp time.Time
}

func (c *Certifications) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*c).Name()) }
func (c *Certifications) GetName() string                               { return c.Name }
func (c *Certifications) SetName(name string)                           { c.Name = name }
func (c *Certifications) GetNamespace() string                          { return c.Namespace }
func (c *Certifications) SetNamespace(namespace string)                 { c.Namespace = namespace }
func (c *Certifications) GetOwnerReference() models.OwnerReference      { return c.OwnerRef }
func (c *Certifications) SetOwnerReference(owner models.OwnerReference) { c.OwnerRef = owner }
func (c *Certifications) GetID() string {
	return strings.ToLower(c.GetKind() + ":" + c.Name + ":" + c.Namespace)
}
func (c *Certifications) GetCreationTimestamp() time.Time          { return c.CreationTimestamp }
func (c *Certifications) SetCreationTimestamp(timestamp time.Time) { c.CreationTimestamp = timestamp }
