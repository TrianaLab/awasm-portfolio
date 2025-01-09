package types

import "awasm-portfolio/internal/models"

type Namespace struct {
	Name      string
	OwnerRefs []models.OwnerReference
}

func (ns *Namespace) GetName() string                                   { return ns.Name }
func (ns *Namespace) SetName(name string)                               { ns.Name = name }
func (ns *Namespace) GetNamespace() string                              { return "" } // Namespaces have no parent namespace
func (ns *Namespace) SetNamespace(namespace string)                     {}
func (ns *Namespace) GetOwnerReferences() []models.OwnerReference       { return ns.OwnerRefs }
func (ns *Namespace) SetOwnerReferences(owners []models.OwnerReference) { ns.OwnerRefs = owners }
