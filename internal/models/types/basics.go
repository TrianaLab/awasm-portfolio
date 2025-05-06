package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Location struct {
	Address     string `json:"address" yaml:"Address,omitempty"`
	PostalCode  string `json:"postalCode" yaml:"PostalCode,omitempty"`
	City        string `json:"city" yaml:"City,omitempty"`
	CountryCode string `json:"countryCode" yaml:"CountryCode,omitempty"`
	Region      string `json:"region" yaml:"Region,omitempty"`
}

type Profile struct {
	Network  string `json:"network" yaml:"Network,omitempty"`
	Username string `json:"username" yaml:"Username,omitempty"`
	Url      string `json:"url" yaml:"URL,omitempty"`
}

type Basics struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	FullName          string                `json:"name" yaml:"FullName,omitempty"`
	Label             string                `json:"label" yaml:"Label,omitempty"`
	Image             string                `json:"image" yaml:"Image,omitempty"`
	Email             string                `json:"email" yaml:"Email,omitempty"`
	Phone             string                `json:"phone" yaml:"Phone,omitempty"`
	Url               string                `json:"url" yaml:"URL,omitempty"`
	Summary           string                `json:"summary" yaml:"Summary,omitempty"`
	Location          Location              `json:"location" yaml:"Location,omitempty"`
	Profiles          []Profile             `json:"profiles" yaml:"Profiles,omitempty"`
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
