package models

import (
	"errors"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var companyTypes []CompanyType

type CompanyType struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Abbreviation string    `json:"abbreviation"`
}

func UnmarshalCompanyTypes() error {
	var cts []map[string]interface{}
	if err := viper.UnmarshalKey("companyTypes", &cts); err != nil {
		return err
	}
	for _, ct := range cts {
		id, err := uuid.Parse(ct["id"].(string))
		if err != nil {
			return err
		}
		companyType := CompanyType{
			ID:           id,
			Name:         ct["name"].(string),
			Abbreviation: ct["abbreviation"].(string),
		}
		companyTypes = append(companyTypes, companyType)
	}
	return nil
}

func ListCompanyTypes(filters map[string]interface{}) ([]CompanyType, error) {
	return companyTypes, nil
}

func GetCompanyType(id uuid.UUID) (*CompanyType, error) {
	for _, companyType := range companyTypes {
		if companyType.ID == id {
			return &companyType, nil
		}
	}
	return nil, errors.New("companyType not found")
}
