package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Course struct {
	Title       string
	Institution string
	Duration    string
}

type Education struct {
	Name              string
	Namespace         string
	OwnerRef          models.OwnerReference
	Courses           []Course
	CreationTimestamp time.Time
}

func (e *Education) GetKind() string                               { return strings.ToLower(reflect.TypeOf(*e).Name()) }
func (e *Education) GetName() string                               { return e.Name }
func (e *Education) SetName(name string)                           { e.Name = name }
func (e *Education) GetNamespace() string                          { return e.Namespace }
func (e *Education) SetNamespace(namespace string)                 { e.Namespace = namespace }
func (e *Education) GetOwnerReference() models.OwnerReference      { return e.OwnerRef }
func (e *Education) SetOwnerReference(owner models.OwnerReference) { e.OwnerRef = owner }
func (e *Education) GetID() string {
	return strings.ToLower(e.GetKind() + ":" + e.Name + ":" + e.Namespace)
}
func (e *Education) GetCreationTimestamp() time.Time          { return e.CreationTimestamp }
func (e *Education) SetCreationTimestamp(timestamp time.Time) { e.CreationTimestamp = timestamp }
