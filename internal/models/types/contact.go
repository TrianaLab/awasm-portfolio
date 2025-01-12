package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Contact struct {
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	Email             string
	LinkedIn          string
	GitHub            string
	CreationTimestamp time.Time
}

func (c *Contact) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*c).Name()) }
func (c *Contact) GetName() string                               { return c.Name }
func (c *Contact) SetName(name string)                           { c.Name = name }
func (c *Contact) GetNamespace() string                          { return c.Namespace }
func (c *Contact) SetNamespace(namespace string)                 { c.Namespace = namespace }
func (c *Contact) GetOwnerReference() models.OwnerReference      { return c.OwnerRef }
func (c *Contact) SetOwnerReference(owner models.OwnerReference) { c.OwnerRef = owner }
func (c *Contact) GetID() string {
	return strings.ToLower(c.GetKind() + ":" + c.Name + ":" + c.Namespace)
}
func (c *Contact) GetCreationTimestamp() time.Time          { return c.CreationTimestamp }
func (c *Contact) SetCreationTimestamp(timestamp time.Time) { c.CreationTimestamp = timestamp }
