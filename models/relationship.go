package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Relationship struct {
	ID           uuid.UUID  `json:"id" sql:"primary key" gorm:"type:uuid"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time ``
	LeftID       uuid.UUID  `json:"left_id" gorm:"type:uuid"`
	LeftCompany  *Company   `json:"left_company" gorm:"foreignkey:LeftID"`
	RightID      uuid.UUID  `json:"right_id" gorm:"type:uuid"`
	RightCompany *Company   `json:"right_company" gorm:"foreignkey:RightID"`
	Quality      string     `json:"quality"`
}

func (r *Relationship) BeforeSave() error {
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

func (r *Relationship) Conflicts() bool {
	return false
}

func (r *Relationship) Validate() error {
	return nil
}

func (r *Relationship) EagerLoad(side int) error {
	switch side {
	case 1:
		c, err := GetCompany(r.LeftID)
		if err != nil {
			return err
		}
		r.LeftCompany = c
	case 2:
		c, err := GetCompany(r.RightID)
		if err != nil {
			return err
		}
		r.RightCompany = c
	}
	return nil
}

func (r *Relationship) Save() error {
	return session.Save(r).Error
}

func (r *Relationship) Delete() error {
	if r.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}
	return session.Delete(r).Error
}

func NewRelationship(r *Relationship) error {
	if r.ID == uuid.Nil {
		var err error
		if r.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if r.Conflicts() {
		return fmt.Errorf("the item already exists")
	}
	if err := session.Create(&r).Error; err != nil {
		return err
	}
	return nil
}

func GetRelationship(id uuid.UUID) (*Relationship, error) {
	var item Relationship
	err := session.Where("id = ?", id).First(&item).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}
