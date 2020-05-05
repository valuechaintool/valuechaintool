package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Sector struct {
	ID        uuid.UUID  `json:"id" sql:"primary key" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time ``
	Name      string     `json:"name"`
}

func (s *Sector) BeforeSave() error {
	if err := s.Validate(); err != nil {
		return err
	}
	return nil
}

func (s *Sector) Conflicts() bool {
	return false
}

func (s *Sector) Validate() error {
	return nil
}

func (s *Sector) Save() error {
	return session.Save(s).Error
}

func (s *Sector) Delete() error {
	if s.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}
	return session.Delete(s).Error
}

func NewSector(s *Sector) error {
	if s.ID == uuid.Nil {
		var err error
		if s.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if s.Conflicts() {
		return fmt.Errorf("the item already exists")
	}
	if err := session.Create(&s).Error; err != nil {
		return err
	}
	return nil
}

func ListSectors(filters map[string]interface{}) ([]Sector, error) {
	var items []Sector
	if err := session.Where(filters).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
