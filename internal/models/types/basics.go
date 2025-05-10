package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Location struct {
	Address     string `json:"address,omitempty" yaml:"Address,omitempty"`
	PostalCode  string `json:"postalCode,omitempty" yaml:"PostalCode,omitempty"`
	City        string `json:"city,omitempty" yaml:"City,omitempty"`
	CountryCode string `json:"countryCode,omitempty" yaml:"CountryCode,omitempty"`
	Region      string `json:"region,omitempty" yaml:"Region,omitempty"`
}

type Profile struct {
	Network  string `json:"network,omitempty" yaml:"Network,omitempty"`
	Username string `json:"username,omitempty" yaml:"Username,omitempty"`
	Url      string `json:"url,omitempty" yaml:"URL,omitempty"`
}

type Basics struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	FullName          string                `json:"name,omitempty" yaml:"FullName,omitempty"`
	Label             string                `json:"label,omitempty" yaml:"Label,omitempty"`
	Image             string                `json:"image,omitempty" yaml:"Image,omitempty"`
	Email             string                `json:"email,omitempty" yaml:"Email,omitempty"`
	Phone             string                `json:"phone,omitempty" yaml:"Phone,omitempty"`
	Url               string                `json:"url,omitempty" yaml:"URL,omitempty"`
	Summary           string                `json:"summary,omitempty" yaml:"Summary,omitempty"`
	Location          Location              `json:"location,omitempty" yaml:"Location,omitempty"`
	Profiles          []Profile             `json:"profiles,omitempty" yaml:"Profiles,omitempty"`
}

func (b *Basics) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*b).Name()) }
func (b *Basics) GetName() string                                { return b.Name }
func (b *Basics) SetName(name string)                            { b.Name = name }
func (b *Basics) GetNamespace() string                           { return b.Namespace }
func (b *Basics) SetNamespace(namespace string)                  { b.Namespace = namespace }
func (b *Basics) GetOwnerReference() models.OwnerReference       { return b.OwnerRef }
func (b *Basics) SetOwnerReference(owners models.OwnerReference) { b.OwnerRef = owners }
func (b *Basics) GetID() string {
	return strings.ToLower(b.GetKind() + ":" + b.Name + ":" + b.Namespace)
}
func (b *Basics) GetCreationTimestamp() time.Time          { return b.CreationTimestamp }
func (b *Basics) SetCreationTimestamp(timestamp time.Time) { b.CreationTimestamp = timestamp }
