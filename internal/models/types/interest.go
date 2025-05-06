package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Interest struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Interest          string                `json:"name" yaml:"Interest,omitempty"`
	Keywords          []string              `json:"keywords" yaml:"Keywords,omitempty"`
}

func (i *Interest) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*i).Name()) }
func (i *Interest) GetName() string                                { return i.Name }
func (i *Interest) SetName(name string)                            { i.Name = name }
func (i *Interest) GetNamespace() string                           { return i.Namespace }
func (i *Interest) SetNamespace(namespace string)                  { i.Namespace = namespace }
func (i *Interest) GetOwnerReference() models.OwnerReference       { return i.OwnerRef }
func (i *Interest) SetOwnerReference(owners models.OwnerReference) { i.OwnerRef = owners }
func (i *Interest) GetID() string {
	return strings.ToLower(i.GetKind() + ":" + i.Name + ":" + i.Namespace)
}
func (i *Interest) GetCreationTimestamp() time.Time          { return i.CreationTimestamp }
func (i *Interest) SetCreationTimestamp(timestamp time.Time) { i.CreationTimestamp = timestamp }
