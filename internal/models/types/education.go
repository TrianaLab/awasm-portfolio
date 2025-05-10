package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Education struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Institution       string                `json:"institution,omitempty" yaml:"Institution,omitempty"`
	URL               string                `json:"url,omitempty" yaml:"URL,omitempty"`
	Area              string                `json:"area,omitempty" yaml:"Area,omitempty"`
	StudyType         string                `json:"studyType,omitempty" yaml:"StudyType,omitempty"`
	StartDate         string                `json:"startDate,omitempty" yaml:"StartDate,omitempty"`
	EndDate           string                `json:"endDate,omitempty" yaml:"EndDate,omitempty"`
	Score             string                `json:"score,omitempty" yaml:"Score,omitempty"`
	Courses           []string              `json:"courses,omitempty" yaml:"Courses,omitempty"`
}

func (e *Education) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*e).Name()) }
func (e *Education) GetName() string                                { return e.Name }
func (e *Education) SetName(name string)                            { e.Name = name }
func (e *Education) GetNamespace() string                           { return e.Namespace }
func (e *Education) SetNamespace(namespace string)                  { e.Namespace = namespace }
func (e *Education) GetOwnerReference() models.OwnerReference       { return e.OwnerRef }
func (e *Education) SetOwnerReference(owners models.OwnerReference) { e.OwnerRef = owners }
func (e *Education) GetID() string {
	return strings.ToLower(e.GetKind() + ":" + e.Name + ":" + e.Namespace)
}
func (e *Education) GetCreationTimestamp() time.Time          { return e.CreationTimestamp }
func (e *Education) SetCreationTimestamp(timestamp time.Time) { e.CreationTimestamp = timestamp }
