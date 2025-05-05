package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Resume struct {
	Kind              string                `json:"kind,omitempty" yaml:"Kind,omitempty"`
	Name              string                `json:"name,omitempty" yaml:"Name,omitempty"`
	Namespace         string                `json:"namespace,omitempty" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"ownerRef,omitempty" yaml:"OwnerRef,omitempty"`
	CreationTimestamp time.Time             `json:"creationTimestamp,omitempty" yaml:"CreationTimestamp,omitempty"`
	Basics            Basics                `json:"basics,omitempty" yaml:"Basics,omitempty"`
	Work              []Work                `json:"work,omitempty" yaml:"Work,omitempty"`
	Volunteer         []Volunteer           `json:"volunteer,omitempty" yaml:"Volunteer,omitempty"`
	Education         []Education           `json:"education,omitempty" yaml:"Education,omitempty"`
	Awards            []Award               `json:"awards,omitempty" yaml:"Awards,omitempty"`
	Certificates      []Certificate         `json:"certificates,omitempty" yaml:"Certificates,omitempty"`
	Publications      []Publication         `json:"publications,omitempty" yaml:"Publications,omitempty"`
	Skills            []Skill               `json:"skills,omitempty" yaml:"Skills,omitempty"`
	Languages         []Language            `json:"languages,omitempty" yaml:"Languages,omitempty"`
	Interests         []Interest            `json:"interests,omitempty" yaml:"Interests,omitempty"`
	References        []Reference           `json:"references,omitempty" yaml:"References,omitempty"`
	Projects          []Project             `json:"projects,omitempty" yaml:"Projects,omitempty"`
}

func (r *Resume) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*r).Name()) }
func (r *Resume) GetName() string                                { return r.Name }
func (r *Resume) SetName(name string)                            { r.Name = name }
func (r *Resume) GetNamespace() string                           { return r.Namespace }
func (r *Resume) SetNamespace(namespace string)                  { r.Namespace = namespace }
func (r *Resume) GetOwnerReference() models.OwnerReference       { return r.OwnerRef }
func (r *Resume) SetOwnerReference(owners models.OwnerReference) { r.OwnerRef = owners }
func (r *Resume) GetID() string {
	return strings.ToLower(r.GetKind() + ":" + r.Name + ":" + r.Namespace)
}
func (r *Resume) GetCreationTimestamp() time.Time          { return r.CreationTimestamp }
func (r *Resume) SetCreationTimestamp(timestamp time.Time) { r.CreationTimestamp = timestamp }
