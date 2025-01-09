package types

import "awasm-portfolio/internal/models"

type Profile struct {
	Name      string
	Namespace string
	OwnerRefs []models.OwnerReference
}

func (p *Profile) GetName() string                                   { return p.Name }
func (p *Profile) SetName(name string)                               { p.Name = name }
func (p *Profile) GetNamespace() string                              { return p.Namespace }
func (p *Profile) SetNamespace(namespace string)                     { p.Namespace = namespace }
func (p *Profile) GetOwnerReferences() []models.OwnerReference       { return p.OwnerRefs }
func (p *Profile) SetOwnerReferences(owners []models.OwnerReference) { p.OwnerRefs = owners }
