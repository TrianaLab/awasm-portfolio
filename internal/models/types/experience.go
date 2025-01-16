package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Job struct {
	Title       string
	Description string
	Company     string
	Duration    string
}

type Experience struct {
	Kind              string
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	Jobs              []Job
	CreationTimestamp time.Time
}

func (e *Experience) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*e).Name()) }
func (e *Experience) GetName() string                               { return e.Name }
func (e *Experience) SetName(name string)                           { e.Name = name }
func (e *Experience) GetNamespace() string                          { return e.Namespace }
func (e *Experience) SetNamespace(namespace string)                 { e.Namespace = namespace }
func (e *Experience) GetOwnerReference() models.OwnerReference      { return e.OwnerRef }
func (e *Experience) SetOwnerReference(owner models.OwnerReference) { e.OwnerRef = owner }
func (e *Experience) GetID() string {
	return strings.ToLower(e.GetKind() + ":" + e.Name + ":" + e.Namespace)
}
func (e *Experience) GetCreationTimestamp() time.Time          { return e.CreationTimestamp }
func (e *Experience) SetCreationTimestamp(timestamp time.Time) { e.CreationTimestamp = timestamp }
