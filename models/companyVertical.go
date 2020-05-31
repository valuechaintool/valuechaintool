package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type CompanyVertical struct {
	CompanyID  uuid.UUID `gorm:"primary_key;auto_increment:false"`
	VerticalID uuid.UUID `gorm:"primary_key;auto_increment:false"`
}

func GetVerticalsByCompany(companyID uuid.UUID) ([]Vertical, error) {
	var verticals []Vertical
	var items []CompanyVertical
	if err := session.Where("company_id = ?", companyID.String()).Find(&items).Error; err != nil {
		return nil, err
	}
	for _, item := range items {
		vertical, err := GetVertical(item.VerticalID)
		if err != nil {
			return nil, err
		}
		verticals = append(verticals, *vertical)
	}
	return verticals, nil
}

func syncCompanyVertical(tx *gorm.DB, companyID uuid.UUID, verticals []Vertical) error {
	if err := tx.Delete(CompanyVertical{}, "company_id = ?", companyID.String()).Error; err != nil {
		return err
	}
	for _, vertical := range verticals {
		cv := CompanyVertical{
			CompanyID:  companyID,
			VerticalID: vertical.ID,
		}
		if err := tx.Create(&cv).Error; err != nil {
			return err
		}
	}
	return nil
}
