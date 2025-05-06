package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Certificate struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Certificate       string                `json:"name" yaml:"Certificate,omitempty"`
	Date              string                `json:"date" yaml:"Date,omitempty"`
	Issuer            string                `json:"issuer" yaml:"Issuer,omitempty"`
	URL               string                `json:"url" yaml:"URL,omitempty"`
}

func (c *Certificate) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*c).Name()) }
func (c *Certificate) GetName() string                                { return c.Name }
func (c *Certificate) SetName(name string)                            { c.Name = name }
func (c *Certificate) GetNamespace() string                           { return c.Namespace }
func (c *Certificate) SetNamespace(namespace string)                  { c.Namespace = namespace }
func (c *Certificate) GetOwnerReference() models.OwnerReference       { return c.OwnerRef }
func (c *Certificate) SetOwnerReference(owners models.OwnerReference) { c.OwnerRef = owners }
func (c *Certificate) GetID() string {
	return strings.ToLower(c.GetKind() + ":" + c.Name + ":" + c.Namespace)
}
func (c *Certificate) GetCreationTimestamp() time.Time          { return c.CreationTimestamp }
func (c *Certificate) SetCreationTimestamp(timestamp time.Time) { c.CreationTimestamp = timestamp }
