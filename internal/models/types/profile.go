package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
)

type Profile struct {
	Name           string
	Namespace      string
	OwnerRef       models.OwnerReference
	Certifications Certifications
	Contact        Contact
	Contributions  Contributions
	Education      Education
	Experience     Experience
	Skills         Skills
}

func (p *Profile) GetKind() string                                { return reflect.TypeOf(*p).Name() }
func (p *Profile) GetName() string                                { return p.Name }
func (p *Profile) SetName(name string)                            { p.Name = name }
func (p *Profile) GetNamespace() string                           { return p.Namespace }
func (p *Profile) SetNamespace(namespace string)                  { p.Namespace = namespace }
func (p *Profile) GetOwnerReference() models.OwnerReference       { return p.OwnerRef }
func (p *Profile) SetOwnerReference(owners models.OwnerReference) { p.OwnerRef = owners }
