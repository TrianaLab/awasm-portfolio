package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Award struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"title" yaml:"Title,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Date              string                `json:"date" yaml:"Date,omitempty"`
	Awarder           string                `json:"awarder" yaml:"Awarder,omitempty"`
	Summary           string                `json:"summary" yaml:"Summary,omitempty"`
}

func (a *Award) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*a).Name()) }
func (a *Award) GetName() string                                { return a.Name }
func (a *Award) SetName(name string)                            { a.Name = name }
func (a *Award) GetNamespace() string                           { return a.Namespace }
func (a *Award) SetNamespace(namespace string)                  { a.Namespace = namespace }
func (a *Award) GetOwnerReference() models.OwnerReference       { return a.OwnerRef }
func (a *Award) SetOwnerReference(owners models.OwnerReference) { a.OwnerRef = owners }
func (a *Award) GetID() string {
	return strings.ToLower(a.GetKind() + ":" + a.Name + ":" + a.Namespace)
}
func (a *Award) GetCreationTimestamp() time.Time          { return a.CreationTimestamp }
func (a *Award) SetCreationTimestamp(timestamp time.Time) { a.CreationTimestamp = timestamp }
