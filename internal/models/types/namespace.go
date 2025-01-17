package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"time"
)

type Namespace struct {
	Kind              string                `json:"Kind,omitempty" yaml:"Kind,omitempty"`
	Name              string                `json:"Name,omitempty" yaml:"Name,omitempty"`
	Namespace         string                `json:"Namespace,omitempty" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"Owner,omitempty" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"CreationTimestamp,omitempty" yaml:"CreationTimestamp,omitempty"`
}

func (ns *Namespace) GetKind() string                               { return "namespace" }
func (ns *Namespace) GetName() string                               { return ns.Name }
func (ns *Namespace) SetName(name string)                           { ns.Name = name }
func (ns *Namespace) GetNamespace() string                          { return "" }
func (ns *Namespace) SetNamespace(namespace string)                 {}
func (ns *Namespace) GetOwnerReference() models.OwnerReference      { return ns.OwnerRef }
func (ns *Namespace) SetOwnerReference(owner models.OwnerReference) {}
func (ns *Namespace) GetID() string {
	return strings.ToLower(ns.GetKind() + ":" + ns.Name + ":" + ns.Namespace)
}
func (ns *Namespace) GetCreationTimestamp() time.Time          { return ns.CreationTimestamp }
func (ns *Namespace) SetCreationTimestamp(timestamp time.Time) { ns.CreationTimestamp = timestamp }
