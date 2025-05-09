package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Reference struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Person            string                `json:"name" yaml:"Person,omitempty"`
	Reference         string                `json:"reference" yaml:"Reference,omitempty"`
}

func (r *Reference) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*r).Name()) }
func (r *Reference) GetName() string                                { return r.Name }
func (r *Reference) SetName(name string)                            { r.Name = name }
func (r *Reference) GetNamespace() string                           { return r.Namespace }
func (r *Reference) SetNamespace(namespace string)                  { r.Namespace = namespace }
func (r *Reference) GetOwnerReference() models.OwnerReference       { return r.OwnerRef }
func (r *Reference) SetOwnerReference(owners models.OwnerReference) { r.OwnerRef = owners }
func (r *Reference) GetID() string {
	return strings.ToLower(r.GetKind() + ":" + r.Name + ":" + r.Namespace)
}
func (r *Reference) GetCreationTimestamp() time.Time          { return r.CreationTimestamp }
func (r *Reference) SetCreationTimestamp(timestamp time.Time) { r.CreationTimestamp = timestamp }
