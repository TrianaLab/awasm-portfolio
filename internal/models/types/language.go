package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Language struct {
	models.Meta `json:"-" yaml:",inline"`
	Language    string `json:"language,omitempty" yaml:"Language,omitempty"`
	Fluency     string `json:"fluency,omitempty" yaml:"Fluency,omitempty"`
}
