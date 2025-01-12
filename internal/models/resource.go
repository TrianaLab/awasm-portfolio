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
	Owner     Resource
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
