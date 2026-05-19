package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Publication struct {
	models.Meta `json:"-" yaml:",inline"`
	Publication string `json:"name,omitempty" yaml:"Publication,omitempty"`
	Publisher   string `json:"publisher,omitempty" yaml:"Publisher,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty" yaml:"ReleaseDate,omitempty"`
	URL         string `json:"url,omitempty" yaml:"URL,omitempty"`
	Summary     string `json:"summary,omitempty" yaml:"Summary,omitempty"`
}
