package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChangeType int

const (
	CompanyCreated ChangeType = iota
	CompanyUpdated
	CompanyDeleted
	VerticalAdded
	VerticalRemoved
	RelationshipCreated
	RelationshipUpdated
	RelationshipDeleted
)

type Change struct {
	ID             uuid.UUID     `json:"id" sql:"primary key" gorm:"type:uuid"`
	CompanyID      uuid.UUID     `json:"-" gorm:"index:changes_company_id"`
	RelationshipID *uuid.UUID    `json:"relationship_id"`
	Relationship   *Relationship `json:"relationship"`
	Time           time.Time     `json:"time"`
	UserID         uuid.UUID     `json:"user_id"`
	User           *User         `json:"user"`
	Type           ChangeType    `json:"type"`
	Key            string        `json:"key"`
	PreviousValue  string        `json:"previous_value"`
	NewValue       string        `json:"new_value"`
}

func newChange(tx *gorm.DB, c *Change) error {
	if c.ID == uuid.Nil {
		var err error
		if c.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	c.Time = time.Now()
	if err := tx.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func companyUpdateChange(tx *gorm.DB, companyID uuid.UUID, userID uuid.UUID, key string, previousValue string, newValue string) error {
	if newValue != previousValue {
		change := Change{
			CompanyID:     companyID,
			UserID:        userID,
			Type:          CompanyUpdated,
			Key:           key,
			PreviousValue: previousValue,
			NewValue:      newValue,
		}
		if err := newChange(tx, &change); err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func ListChangesByCompany(companyID uuid.UUID, eagerLoad bool) ([]Change, error) {
	var items []Change
	var err error
	if err := session.Where("company_id = ?", companyID).Order("time desc").Find(&items).Error; err != nil {
		return nil, err
	}
	for c := range items {
		if eagerLoad {
			if items[c].RelationshipID != nil {
				if items[c].Relationship, err = GetRelationshipWithDeleted(*items[c].RelationshipID); err != nil {
					return nil, err
				}
				if items[c].Relationship.LeftID != companyID {
					rel := items[c].Relationship.Reverse()
					items[c].Relationship = &rel
				}
				if err := items[c].Relationship.EagerLoad(2); err != nil {
					return nil, err
				}
			}
			if items[c].User, err = GetUser(items[c].UserID); err != nil {
				return nil, err
			}
		}
	}
	return items, nil
}
