package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Profile struct {
	Kind              string
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	Certifications    Certifications
	Contact           Contact
	Contributions     Contributions
	Education         Education
	Experience        Experience
	Skills            Skills
	CreationTimestamp time.Time
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
