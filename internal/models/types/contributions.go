package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"time"
)

type Contribution struct {
	Project     string `yaml:"Project,omitempty"`
	Description string `yaml:"Description,omitempty"`
	Link        string `yaml:"Link,omitempty"`
}

type Contributions struct {
	Kind              string                `yaml:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty"`
	Contributions     []Contribution        `yaml:"Contributions,omitempty"`
}

func (c *Contributions) GetKind() string                               { return "contributions" }
func (c *Contributions) GetName() string                               { return c.Name }
func (c *Contributions) SetName(name string)                           { c.Name = name }
func (c *Contributions) GetNamespace() string                          { return c.Namespace }
func (c *Contributions) SetNamespace(namespace string)                 { c.Namespace = namespace }
func (c *Contributions) GetOwnerReference() models.OwnerReference      { return c.OwnerRef }
func (c *Contributions) SetOwnerReference(owner models.OwnerReference) { c.OwnerRef = owner }
func (c *Contributions) GetID() string {
	return strings.ToLower(c.GetKind() + ":" + c.Name + ":" + c.Namespace)
}
func (c *Contributions) GetCreationTimestamp() time.Time          { return c.CreationTimestamp }
func (c *Contributions) SetCreationTimestamp(timestamp time.Time) { c.CreationTimestamp = timestamp }
