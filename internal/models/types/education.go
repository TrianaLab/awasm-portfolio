package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Education struct {
	models.Meta `json:"-" yaml:",inline"`
	Institution string   `json:"institution,omitempty" yaml:"Institution,omitempty"`
	URL         string   `json:"url,omitempty" yaml:"URL,omitempty"`
	Area        string   `json:"area,omitempty" yaml:"Area,omitempty"`
	StudyType   string   `json:"studyType,omitempty" yaml:"StudyType,omitempty"`
	StartDate   string   `json:"startDate,omitempty" yaml:"StartDate,omitempty"`
	EndDate     string   `json:"endDate,omitempty" yaml:"EndDate,omitempty"`
	Score       string   `json:"score,omitempty" yaml:"Score,omitempty"`
	Courses     []string `json:"courses,omitempty" yaml:"Courses,omitempty"`
}
