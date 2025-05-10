package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Language struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Language          string                `json:"language,omitempty" yaml:"Language,omitempty"`
	Fluency           string                `json:"fluency,omitempty" yaml:"Fluency,omitempty"`
}

func (l *Language) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*l).Name()) }
func (l *Language) GetName() string                                { return l.Name }
func (l *Language) SetName(name string)                            { l.Name = name }
func (l *Language) GetNamespace() string                           { return l.Namespace }
func (l *Language) SetNamespace(namespace string)                  { l.Namespace = namespace }
func (l *Language) GetOwnerReference() models.OwnerReference       { return l.OwnerRef }
func (l *Language) SetOwnerReference(owners models.OwnerReference) { l.OwnerRef = owners }
func (l *Language) GetID() string {
	return strings.ToLower(l.GetKind() + ":" + l.Name + ":" + l.Namespace)
}
func (l *Language) GetCreationTimestamp() time.Time          { return l.CreationTimestamp }
func (l *Language) SetCreationTimestamp(timestamp time.Time) { l.CreationTimestamp = timestamp }
