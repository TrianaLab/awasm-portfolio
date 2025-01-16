package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Profile struct {
	Kind              string                `yaml:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty"`
	Contributions     Contributions         `yaml:"Contributions,omitempty"`
	Experience        Experience            `yaml:"Experience,omitempty"`
	Certifications    Certifications        `yaml:"Certifications,omitempty"`
	Education         Education             `yaml:"Education,omitempty"`
	Skills            Skills                `yaml:"Skills,omitempty"`
	Contact           Contact               `yaml:"Contact,omitempty"`
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
