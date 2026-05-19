package types

import "github.com/TrianaLab/awasm-portfolio/internal/models"

type Location struct {
	Address     string `json:"address,omitempty" yaml:"Address,omitempty"`
	PostalCode  string `json:"postalCode,omitempty" yaml:"PostalCode,omitempty"`
	City        string `json:"city,omitempty" yaml:"City,omitempty"`
	CountryCode string `json:"countryCode,omitempty" yaml:"CountryCode,omitempty"`
	Region      string `json:"region,omitempty" yaml:"Region,omitempty"`
}

type Profile struct {
	Network  string `json:"network,omitempty" yaml:"Network,omitempty"`
	Username string `json:"username,omitempty" yaml:"Username,omitempty"`
	Url      string `json:"url,omitempty" yaml:"URL,omitempty"`
}

type Basics struct {
	models.Meta `json:"-" yaml:",inline"`
	FullName    string    `json:"name,omitempty" yaml:"FullName,omitempty"`
	Label       string    `json:"label,omitempty" yaml:"Label,omitempty"`
	Image       string    `json:"image,omitempty" yaml:"Image,omitempty"`
	Email       string    `json:"email,omitempty" yaml:"Email,omitempty"`
	Phone       string    `json:"phone,omitempty" yaml:"Phone,omitempty"`
	Url         string    `json:"url,omitempty" yaml:"URL,omitempty"`
	Summary     string    `json:"summary,omitempty" yaml:"Summary,omitempty"`
	Location    Location  `json:"location,omitempty" yaml:"Location,omitempty"`
	Profiles    []Profile `json:"profiles,omitempty" yaml:"Profiles,omitempty"`
}
