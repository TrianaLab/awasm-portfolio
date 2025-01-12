package models

import "strings"

type Resource interface {
	GetKind() string
	GetName() string
	SetName(name string)
	GetNamespace() string
	SetNamespace(namespace string)
	GetOwnerReference() OwnerReference
	SetOwnerReference(owner OwnerReference)
	GetID() string
}

type OwnerReference struct {
	Kind      string
	Name      string
	Namespace string
}

func (o OwnerReference) GetID() string {
	return strings.ToLower(o.Kind + ":" + o.Name + ":" + o.Namespace)
}
