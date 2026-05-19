package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Skill struct {
	models.Meta `json:"-" yaml:",inline"`
	Skill       string   `json:"name,omitempty" yaml:"Skill,omitempty"`
	Level       string   `json:"level,omitempty" yaml:"Level,omitempty"`
	Keywords    []string `json:"keywords,omitempty" yaml:"Keywords,omitempty"`
}
