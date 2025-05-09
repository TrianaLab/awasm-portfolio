package types

import (
	"awasm-portfolio/internal/models"
	"reflect"
	"strings"
	"time"
)

type Work struct {
	Kind              string                `json:"-" yaml:"Kind,omitempty"`
	Name              string                `json:"-" yaml:"Name,omitempty"`
	Namespace         string                `json:"-" yaml:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `json:"-" yaml:"Owner,omitempty"`
	CreationTimestamp time.Time             `json:"-" yaml:"CreationTimestamp,omitempty"`
	Company           string                `json:"name" yaml:"Company,omitempty"`
	Position          string                `json:"position" yaml:"Position,omitempty"`
	URL               string                `json:"url" yaml:"URL,omitempty"`
	StartDate         string                `json:"startDate" yaml:"StartDate,omitempty"`
	EndDate           string                `json:"endDate" yaml:"EndDate,omitempty"`
	Summary           string                `json:"summary" yaml:"Summary,omitempty"`
	Highlights        []string              `json:"highlights" yaml:"Highlights,omitempty"`
}

func (w *Work) GetKind() string                                { return strings.ToLower(reflect.TypeOf(*w).Name()) }
func (w *Work) GetName() string                                { return w.Name }
func (w *Work) SetName(name string)                            { w.Name = name }
func (w *Work) GetNamespace() string                           { return w.Namespace }
func (w *Work) SetNamespace(namespace string)                  { w.Namespace = namespace }
func (w *Work) GetOwnerReference() models.OwnerReference       { return w.OwnerRef }
func (w *Work) SetOwnerReference(owners models.OwnerReference) { w.OwnerRef = owners }
func (w *Work) GetID() string {
	return strings.ToLower(w.GetKind() + ":" + w.Name + ":" + w.Namespace)
}
func (w *Work) GetCreationTimestamp() time.Time          { return w.CreationTimestamp }
func (w *Work) SetCreationTimestamp(timestamp time.Time) { w.CreationTimestamp = timestamp }
