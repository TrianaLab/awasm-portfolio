package models

type Resource interface {
	GetName() string
	SetName(name string)
	GetNamespace() string
	SetNamespace(namespace string)
	GetOwnerReferences() []OwnerReference
	SetOwnerReferences(owners []OwnerReference)
}

type OwnerReference struct {
	Kind string
	Name string
}
