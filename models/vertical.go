package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var verticals []Vertical

type Vertical struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func UnmarshalVerticals() error {
	var ss []map[string]interface{}
	if err := viper.UnmarshalKey("verticals", &ss); err != nil {
		return err
	}
	for _, s := range ss {
		id, err := uuid.Parse(s["id"].(string))
		if err != nil {
			return err
		}
		vertical := Vertical{
			ID:   id,
			Name: s["name"].(string),
		}
		verticals = append(verticals, vertical)
	}
	return nil
}

func ListVerticals(filters map[string]interface{}) ([]Vertical, error) {
	return verticals, nil
}

func GetVertical(id uuid.UUID) (*Vertical, error) {
	for _, vertical := range verticals {
		if vertical.ID == id {
			return &vertical, nil
		}
	}
	return nil, errors.New("vertical not found")
}

func IsInVerticalArray(v Vertical, verticals []Vertical) bool {
	for _, vertical := range verticals {
		if v.ID == vertical.ID {
			return true
		}
	}
	return false
}
