package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Publication struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Publication       string                `json:"name" yaml:"Publication,omitempty"`
	Publisher         string                `json:"publisher" yaml:"Publisher,omitempty"`
	ReleaseDate       string                `json:"releaseDate" yaml:"ReleaseDate,omitempty"`
	URL               string                `json:"url" yaml:"URL,omitempty"`
	Summary           string                `json:"summary" yaml:"Summary,omitempty"`
}

func (p *Publication) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*p).Name()) }
func (p *Publication) GetName() string                                { return p.Name }
func (p *Publication) SetName(name string)                            { p.Name = name }
func (p *Publication) GetNamespace() string                           { return p.Namespace }
func (p *Publication) SetNamespace(namespace string)                  { p.Namespace = namespace }
func (p *Publication) GetOwnerReference() models.OwnerReference       { return p.OwnerRef }
func (p *Publication) SetOwnerReference(owners models.OwnerReference) { p.OwnerRef = owners }
func (p *Publication) GetID() string {
	return strings.ToLower(p.GetKind() + ":" + p.Name + ":" + p.Namespace)
}
func (p *Publication) GetCreationTimestamp() time.Time          { return p.CreationTimestamp }
func (p *Publication) SetCreationTimestamp(timestamp time.Time) { p.CreationTimestamp = timestamp }
