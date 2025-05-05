package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Education struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"institution" yaml:"Institution,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	URL               string                `json:"url" yaml:"URL,omitempty"`
	Area              string                `json:"area" yaml:"Area,omitempty"`
	StudyType         string                `json:"studyType" yaml:"StudyType,omitempty"`
	StartDate         string                `json:"startDate" yaml:"StartDate,omitempty"`
	EndDate           string                `json:"endDate" yaml:"EndDate,omitempty"`
	Score             string                `json:"score" yaml:"Score,omitempty"`
	Courses           []string              `json:"courses" yaml:"Courses,omitempty"`
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
