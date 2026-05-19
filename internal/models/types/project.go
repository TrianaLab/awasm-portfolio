package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Project struct {
	models.Meta `json:"-" yaml:",inline"`
	Project     string   `json:"name,omitempty" yaml:"Project,omitempty"`
	StartDate   string   `json:"startDate,omitempty" yaml:"StartDate,omitempty"`
	EndDate     string   `json:"endDate,omitempty" yaml:"EndDate,omitempty"`
	Description string   `json:"description,omitempty" yaml:"Description,omitempty"`
	Highlights  []string `json:"highlights,omitempty" yaml:"Highlights,omitempty"`
	URL         string   `json:"url,omitempty" yaml:"URL,omitempty"`
}
