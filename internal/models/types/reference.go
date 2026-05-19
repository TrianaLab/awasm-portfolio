package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Reference struct {
	models.Meta `json:"-" yaml:",inline"`
	Person      string `json:"name,omitempty" yaml:"Person,omitempty"`
	Reference   string `json:"reference,omitempty" yaml:"Reference,omitempty"`
}
