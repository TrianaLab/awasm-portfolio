package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Certificate struct {
	models.Meta `json:"-" yaml:",inline"`
	Certificate string `json:"name,omitempty" yaml:"Certificate,omitempty"`
	Date        string `json:"date,omitempty" yaml:"Date,omitempty"`
	Issuer      string `json:"issuer,omitempty" yaml:"Issuer,omitempty"`
	URL         string `json:"url,omitempty" yaml:"URL,omitempty"`
}
