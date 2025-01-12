package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
)

type Course struct {
	Title       string
	Institution string
	Duration    string
}

type Education struct {
	Name      string
	Namespace string
	OwnerRef  models.OwnerReference
	Courses   []Course
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
