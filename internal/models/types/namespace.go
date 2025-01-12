package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Namespace struct {
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	CreationTimestamp time.Time
}

func (ns *Namespace) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*ns).Name()) }
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
