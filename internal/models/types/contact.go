package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"time"
)

type Contact struct {
	Kind              string                `yaml:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty"`
	Email             string                `yaml:"Email,omitempty"`
	Linkedin          string                `yaml:"LinkedIn,omitempty"`
	Github            string                `yaml:"GitHub,omitempty"`
}

func (c *Contact) GetKind() string                               { return "contact" }
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
