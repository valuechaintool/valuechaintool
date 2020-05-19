package models

type Tier struct {
	ID    int    `mapstructure:"id"`
	Name  string `mapstructure:"name"`
	Color string `mapstructure: "color"`
}
