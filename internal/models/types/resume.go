package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Resume struct {
	models.Meta  `json:"-" yaml:",inline"`
	Basics       Basics        `json:"basics,omitempty" yaml:"Basics,omitempty"`
	Work         []Work        `json:"work,omitempty" yaml:"Work,omitempty"`
	Volunteer    []Volunteer   `json:"volunteer,omitempty" yaml:"Volunteer,omitempty"`
	Education    []Education   `json:"education,omitempty" yaml:"Education,omitempty"`
	Awards       []Award       `json:"awards,omitempty" yaml:"Awards,omitempty"`
	Certificates []Certificate `json:"certificates,omitempty" yaml:"Certificates,omitempty"`
	Publications []Publication `json:"publications,omitempty" yaml:"Publications,omitempty"`
	Skills       []Skill       `json:"skills,omitempty" yaml:"Skills,omitempty"`
	Languages    []Language    `json:"languages,omitempty" yaml:"Languages,omitempty"`
	Interests    []Interest    `json:"interests,omitempty" yaml:"Interests,omitempty"`
	References   []Reference   `json:"references,omitempty" yaml:"References,omitempty"`
	Projects     []Project     `json:"projects,omitempty" yaml:"Projects,omitempty"`
}
