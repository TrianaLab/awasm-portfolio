package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Profile struct {
	Kind              string                `json:"Kind,omitempty" yaml:"Kind,omitempty"`
	Name              string                `json:"Name,omitempty" yaml:"Name,omitempty"`
	Namespace         string                `json:"Namespace,omitempty" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"CreationTimestamp,omitempty" yaml:"CreationTimestamp,omitempty"`
	Contributions     Contributions         `json:"Contributions,omitempty" yaml:"Contributions,omitempty"`
	Experience        Experience            `json:"Experience,omitempty" yaml:"Experience,omitempty"`
	Certifications    Certifications        `json:"Certifications,omitempty" yaml:"Certifications,omitempty"`
	Education         Education             `json:"Education,omitempty" yaml:"Education,omitempty"`
	Skills            Skills                `json:"Skills,omitempty" yaml:"Skills,omitempty"`
	Contact           Contact               `json:"Contact,omitempty" yaml:"Contact,omitempty"`
}

func (p *Profile) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*p).Name()) }
func (p *Profile) GetName() string                                { return p.Name }
func (p *Profile) SetName(name string)                            { p.Name = name }
func (p *Profile) GetNamespace() string                           { return p.Namespace }
func (p *Profile) SetNamespace(namespace string)                  { p.Namespace = namespace }
func (p *Profile) GetOwnerReference() models.OwnerReference       { return p.OwnerRef }
func (p *Profile) SetOwnerReference(owners models.OwnerReference) { p.OwnerRef = owners }
func (p *Profile) GetID() string {
	return strings.ToLower(p.GetKind() + ":" + p.Name + ":" + p.Namespace)
}
func (p *Profile) GetCreationTimestamp() time.Time          { return p.CreationTimestamp }
func (p *Profile) SetCreationTimestamp(timestamp time.Time) { p.CreationTimestamp = timestamp }
