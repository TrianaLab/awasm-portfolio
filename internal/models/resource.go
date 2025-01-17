package models

import (
	"strings"
	"time"
)

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
