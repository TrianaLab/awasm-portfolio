package types

import (
	"awasm-portfolio/internal/models"
	"strings"
	"time"
)

type Job struct {
	Title       string `yaml:"Title,omitempty" json:"Title,omitempty"`
	Description string `yaml:"Description,omitempty" json:"Description,omitempty"`
	Company     string `yaml:"Company,omitempty" json:"Company,omitempty"`
	Duration    string `yaml:"Duration,omitempty" json:"Duration,omitempty"`
}

type Experience struct {
	Kind              string                `yaml:"Kind,omitempty" json:"Kind,omitempty"`
	Name              string                `yaml:"Name,omitempty" json:"Name,omitempty"`
	Namespace         string                `yaml:"Namespace,omitempty" json:"Namespace,omitempty"`
	OwnerRef          models.OwnerReference `yaml:"Owner,omitempty" json:"Owner,omitempty"`
	CreationTimestamp time.Time             `yaml:"CreationTimestamp,omitempty" json:"CreationTimestamp,omitempty"`
	Jobs              []Job                 `yaml:"Jobs,omitempty" json:"Jobs,omitempty"`
}

func (e *Experience) GetKind() string                               { return "experience" }
func (e *Experience) GetName() string                               { return e.Name }
func (e *Experience) SetName(name string)                           { e.Name = name }
func (e *Experience) GetNamespace() string                          { return e.Namespace }
func (e *Experience) SetNamespace(namespace string)                 { e.Namespace = namespace }
func (e *Experience) GetOwnerReference() models.OwnerReference      { return e.OwnerRef }
func (e *Experience) SetOwnerReference(owner models.OwnerReference) { e.OwnerRef = owner }
func (e *Experience) GetID() string {
	return strings.ToLower(e.GetKind() + ":" + e.Name + ":" + e.Namespace)
}
func (e *Experience) GetCreationTimestamp() time.Time          { return e.CreationTimestamp }
func (e *Experience) SetCreationTimestamp(timestamp time.Time) { e.CreationTimestamp = timestamp }
