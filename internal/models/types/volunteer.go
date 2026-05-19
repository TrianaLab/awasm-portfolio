package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Volunteer struct {
	models.Meta  `json:"-" yaml:",inline"`
	Organization string   `json:"organization,omitempty" yaml:"Organization,omitempty"`
	Position     string   `json:"position,omitempty" yaml:"Position,omitempty"`
	URL          string   `json:"url,omitempty" yaml:"URL,omitempty"`
	StartDate    string   `json:"startDate,omitempty" yaml:"StartDate,omitempty"`
	EndDate      string   `json:"endDate,omitempty" yaml:"EndDate,omitempty"`
	Summary      string   `json:"summary,omitempty" yaml:"Summary,omitempty"`
	Highlights   []string `json:"highlights,omitempty" yaml:"Highlights,omitempty"`
}
