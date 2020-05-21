package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var sectors []Sector

type Sector struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func UnmarshalSectors() error {
	var ss []map[string]interface{}
	if err := viper.UnmarshalKey("sectors", &ss); err != nil {
		return err
	}
	for _, s := range ss {
		id, err := uuid.Parse(s["id"].(string))
		if err != nil {
			return err
		}
		sector := Sector{
			ID:   id,
			Name: s["name"].(string),
		}
		sectors = append(sectors, sector)
	}
	return nil
}

func ListSectors(filters map[string]interface{}) ([]Sector, error) {
	return sectors, nil
}

func GetSector(id uuid.UUID) (*Sector, error) {
	for _, sector := range sectors {
		if sector.ID == id {
			return &sector, nil
		}
	}
	return nil, errors.New("sector not found")
}
