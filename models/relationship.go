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
	LeftTier     int        `json:"left_tier"`
	RightTier    int        `json:"right_tier"`
	Notes        string     `json:"notes"`
}

func (r *Relationship) BeforeSave() error {
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

func (r *Relationship) Reverse() Relationship {
	return Relationship{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
		LeftID:    r.RightID,
		RightID:   r.LeftID,
		LeftTier:  r.RightTier,
		RightTier: r.LeftTier,
		Notes:     r.Notes,
	}
}

func (r *Relationship) Conflicts() bool {
	return false
}

func (r *Relationship) Update(relTier int, notes string) error {
	or, err := GetRelationship(r.ID)
	if err != nil {
		return err
	}
	if or.LeftID == r.LeftID {
		session.Model(&or).Update("left_tier", relTier)
	} else {
		session.Model(&or).Update("right_tier", relTier)
	}
	session.Model(&or).Update("notes", notes)
	nr, err := GetRelationship(r.ID)
	if err != nil {
		return err
	}
	if nr.LeftID != r.LeftID {
		*nr = nr.Reverse()
	}
	r = nr //nolint:staticcheck
	return nil
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

func GetRelationshipsByMember(id uuid.UUID, eagerLoad bool) ([]Relationship, error) {
	var items []Relationship
	if err := session.Where("left_id = ?", id).Find(&items).Error; err != nil {
		return nil, err
	}
	if eagerLoad {
		for i := range items {
			if err := items[i].EagerLoad(2); err != nil {
				return nil, err
			}
		}
	}
	var rightItems []Relationship
	if err := session.Where("right_id = ?", id).Find(&rightItems).Error; err != nil {
		return nil, err
	}
	if eagerLoad {
		for _, r := range rightItems {
			item := r.Reverse()
			if err := item.EagerLoad(2); err != nil {
				return nil, err
			}
			items = append(items, item)
		}
	} else {
		items = append(items, rightItems...)
	}
	return items, nil
}
