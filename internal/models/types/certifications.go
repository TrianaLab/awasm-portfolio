package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"time"
)

type Certification struct {
	Description string `yaml:"Description,omitempty"`
	Link        string `yaml:"Link,omitempty"`
}

type Certifications struct {
	Kind              string                `yaml:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty"`
	Certifications    []Certification       `yaml:"Certifications,omitempty"`
}

func (c *Certifications) GetKind() string                               { return "certifications" }
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
