package models

var tiers []Tier

type Tier struct {
	ID    int    `mapstructure:"id" json:"id"`
	Name  string `mapstructure:"name" json:"name"`
	Color string `mapstructure: "color" json:"color"`
}

func ListTiers() []Tier {
	return tiers
}
