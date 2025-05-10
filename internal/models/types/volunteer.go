package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Volunteer struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Organization      string                `json:"organization,omitempty" yaml:"Organization,omitempty"`
	Position          string                `json:"position,omitempty" yaml:"Position,omitempty"`
	URL               string                `json:"url,omitempty" yaml:"URL,omitempty"`
	StartDate         string                `json:"startDate,omitempty" yaml:"StartDate,omitempty"`
	EndDate           string                `json:"endDate,omitempty" yaml:"EndDate,omitempty"`
	Summary           string                `json:"summary,omitempty" yaml:"Summary,omitempty"`
	Highlights        []string              `json:"highlights,omitempty" yaml:"Highlights,omitempty"`
}

func (v *Volunteer) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*v).Name()) }
func (v *Volunteer) GetName() string                                { return v.Name }
func (v *Volunteer) SetName(name string)                            { v.Name = name }
func (v *Volunteer) GetNamespace() string                           { return v.Namespace }
func (v *Volunteer) SetNamespace(namespace string)                  { v.Namespace = namespace }
func (v *Volunteer) GetOwnerReference() models.OwnerReference       { return v.OwnerRef }
func (v *Volunteer) SetOwnerReference(owners models.OwnerReference) { v.OwnerRef = owners }
func (v *Volunteer) GetID() string {
	return strings.ToLower(v.GetKind() + ":" + v.Name + ":" + v.Namespace)
}
func (v *Volunteer) GetCreationTimestamp() time.Time          { return v.CreationTimestamp }
func (v *Volunteer) SetCreationTimestamp(timestamp time.Time) { v.CreationTimestamp = timestamp }
