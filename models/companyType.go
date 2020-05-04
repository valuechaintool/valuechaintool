package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CompanyType struct {
	ID           uuid.UUID  `json:"id" sql:"primary key" gorm:"type:uuid"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time ``
	Name         string     `json:"name"`
	Abbreviation string     `json:"abbreviation"`
}

func (t *CompanyType) BeforeSave() error {
	if err := t.Validate(); err != nil {
		return err
	}
	return nil
}

func (t *CompanyType) Conflicts() bool {
	return false
}

func (t *CompanyType) Validate() error {
	return nil
}

func (t *CompanyType) Save() error {
	return session.Save(t).Error
}

func (t *CompanyType) Delete() error {
	if t.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}
	return session.Delete(t).Error
}

func NewCompanyType(t *CompanyType) error {
	if t.ID == uuid.Nil {
		var err error
		if t.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if t.Conflicts() {
		return fmt.Errorf("the item already exists")
	}
	if err := session.Create(&t).Error; err != nil {
		return err
	}
	return nil
}

func ListCompanyTypes(filters map[string]interface{}) ([]CompanyType, error) {
	var items []CompanyType
	if err := session.Where(filters).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
