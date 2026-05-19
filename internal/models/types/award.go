package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Award struct {
	models.Meta `json:"-" yaml:",inline"`
	Title       string `json:"title,omitempty" yaml:"Title,omitempty"`
	Date        string `json:"date,omitempty" yaml:"Date,omitempty"`
	Awarder     string `json:"awarder,omitempty" yaml:"Awarder,omitempty"`
	Summary     string `json:"summary,omitempty" yaml:"Summary,omitempty"`
}
