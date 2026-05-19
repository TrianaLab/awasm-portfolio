package models

import (
	"strings"
	"time"
)

// Resource is the contract every kubectl-style portfolio object implements.
// All concrete types satisfy it by embedding Meta.
type Resource interface {
	GetKind() string
	GetName() string
	SetName(name string)
	GetNamespace() string
	SetNamespace(namespace string)
	GetOwnerReference() OwnerReference
	SetOwnerReference(owner OwnerReference)
	GetID() string
	GetCreationTimestamp() time.Time
	SetCreationTimestamp(timestamp time.Time)
}

// Meta carries the Kubernetes-style identity + ownership fields shared by
// every Resource. Concrete types embed Meta so they inherit the full
// Resource interface for free.
type Meta struct {
	Kind              string         `json:"-" yaml:"Kind,omitempty"`
	Name              string         `json:"-" yaml:"Name,omitempty"`
	Namespace         string         `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time      `json:"-" yaml:"CreationTimestamp,omitempty"`
}

func (m *Meta) GetKind() string                          { return m.Kind }
func (m *Meta) GetName() string                          { return m.Name }
func (m *Meta) SetName(name string)                      { m.Name = name }
func (m *Meta) GetNamespace() string                     { return m.Namespace }
func (m *Meta) SetNamespace(namespace string)            { m.Namespace = namespace }
func (m *Meta) GetOwnerReference() OwnerReference        { return m.OwnerRef }
func (m *Meta) SetOwnerReference(owner OwnerReference)   { m.OwnerRef = owner }
func (m *Meta) GetCreationTimestamp() time.Time          { return m.CreationTimestamp }
func (m *Meta) SetCreationTimestamp(timestamp time.Time) { m.CreationTimestamp = timestamp }
func (m *Meta) GetID() string {
	return strings.ToLower(m.Kind + ":" + m.Name + ":" + m.Namespace)
}

// OwnerReference points at the resource that owns this one. The optional
// Owner field is a back-pointer kept out of serialization.
type OwnerReference struct {
	Kind      string   `json:"Kind,omitempty" yaml:"Kind,omitempty"`
	Name      string   `json:"Name,omitempty" yaml:"Name,omitempty"`
	Namespace string   `json:"Namespace,omitempty" yaml:"Namespace,omitempty"`
	Owner     Resource `json:"-" yaml:"-"`
}

func (o OwnerReference) GetID() string {
	return strings.ToLower(o.Kind + ":" + o.Name + ":" + o.Namespace)
}

func (o *OwnerReference) GetName() string {
	if o.Owner != nil {
		return o.Owner.GetName()
	}
	return o.Name
}
