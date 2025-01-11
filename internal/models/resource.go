package models

type Resource interface {
	GetKind() string
	GetName() string
	SetName(name string)
	GetNamespace() string
	SetNamespace(namespace string)
	GetOwnerReference() OwnerReference
	SetOwnerReference(owner OwnerReference)
}

type OwnerReference struct {
	Kind string
	Name string
}
