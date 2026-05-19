package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Interest struct {
	models.Meta `json:"-" yaml:",inline"`
	Interest    string   `json:"name,omitempty" yaml:"Interest,omitempty"`
	Keywords    []string `json:"keywords,omitempty" yaml:"Keywords,omitempty"`
}
