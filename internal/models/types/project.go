package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Project struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Project           string                `json:"name,omitempty" yaml:"Project,omitempty"`
	StartDate         string                `json:"startDate,omitempty" yaml:"StartDate,omitempty"`
	EndDate           string                `json:"endDate,omitempty" yaml:"EndDate,omitempty"`
	Description       string                `json:"description,omitempty" yaml:"Description,omitempty"`
	Highlights        []string              `json:"highlights,omitempty" yaml:"Highlights,omitempty"`
	URL               string                `json:"url,omitempty" yaml:"URL,omitempty"`
}

func (p *Project) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*p).Name()) }
func (p *Project) GetName() string                                { return p.Name }
func (p *Project) SetName(name string)                            { p.Name = name }
func (p *Project) GetNamespace() string                           { return p.Namespace }
func (p *Project) SetNamespace(namespace string)                  { p.Namespace = namespace }
func (p *Project) GetOwnerReference() models.OwnerReference       { return p.OwnerRef }
func (p *Project) SetOwnerReference(owners models.OwnerReference) { p.OwnerRef = owners }
func (p *Project) GetID() string {
	return strings.ToLower(p.GetKind() + ":" + p.Name + ":" + p.Namespace)
}
func (p *Project) GetCreationTimestamp() time.Time          { return p.CreationTimestamp }
func (p *Project) SetCreationTimestamp(timestamp time.Time) { p.CreationTimestamp = timestamp }
