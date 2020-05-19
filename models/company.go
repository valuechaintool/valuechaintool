package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Company struct {
	ID            uuid.UUID      `json:"id" sql:"primary key" gorm:"type:uuid"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     *time.Time     ``
	Name          string         `json:"name"`
	TypeID        uuid.UUID      `gorm:"type:uuid"`
	Type          CompanyType    `json:"type" gorm:"foreignkey:TypeID"`
	SectorID      uuid.UUID      `gorm:"type:uuid"`
	Sector        Sector         `json:"sector" gorm:"foreignkey:SectorID"`
	Country       string         `json:"country"`
	Relationships []Relationship `json:"relationships" gorm:"-"`
}

func (c *Company) BeforeSave() error {
	if err := c.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *Company) Conflicts() bool {
	return false
}

func (c *Company) Validate() error {
	return nil
}

func (c *Company) EagerLoad() error {
	var err error
	c.Relationships, err = ListRelationshipsByMember(c.ID)
	return err
}

func (c *Company) Save() error {
	return session.Save(c).Error
}

func (c *Company) Delete() error {
	if c.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}
	return session.Delete(c).Error
}

func NewCompany(c *Company) error {
	if c.ID == uuid.Nil {
		var err error
		if c.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if c.Conflicts() {
		return fmt.Errorf("the item already exists")
	}
	if err := session.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func ListCompanies(filters map[string]interface{}) ([]Company, error) {
	var items []Company
	if err := session.Where(filters).Find(&items).Error; err != nil {
		return nil, err
	}
	for c := range items {
		session.Model(items[c]).Related(&items[c].Type, "Type")
		session.Model(items[c]).Related(&items[c].Sector, "Sector")
	}
	return items, nil
}

func SearchCompanies(query string) ([]Company, error) {
	var items []Company
	if err := session.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(query)+"%").Find(&items).Error; err != nil {
		return nil, err
	}
	for c := range items {
		session.Model(items[c]).Related(&items[c].Type, "Type")
		session.Model(items[c]).Related(&items[c].Sector, "Sector")
	}
	return items, nil
}

func GetCompany(id uuid.UUID) (*Company, error) {
	var item Company
	err := session.Where("id = ?", id).First(&item).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	session.Model(item).Related(&item.Type, "Type")
	session.Model(item).Related(&item.Sector, "Sector")
	return &item, nil
}
