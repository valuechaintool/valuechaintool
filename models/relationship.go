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
	var count int
	if err := session.Model(&Relationship{}).Where("left_id = ?", r.LeftID).Where("right_id = ?", r.RightID).Count(&count).Error; err != nil {
		return true
	}
	if count > 0 {
		return true
	}
	if err := session.Model(&Relationship{}).Where("left_id = ?", r.RightID).Where("right_id = ?", r.LeftID).Count(&count).Error; err != nil {
		return true
	}
	if count > 0 {
		return true
	}
	return false
}

func (r *Relationship) Update(relTier int, notes string, userID uuid.UUID) error {
	or, err := GetRelationship(r.ID)
	if err != nil {
		return err
	}

	tx := session.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	if or.LeftID == r.LeftID {
		if or.LeftTier != relTier {
			change := Change{
				UserID:         userID,
				CompanyID:      or.LeftID,
				RelationshipID: &or.ID,
				Type:           RelationshipUpdated,
				Key:            "tier",
				PreviousValue:  tiers[or.LeftTier].Name,
				NewValue:       tiers[relTier].Name,
			}
			if err := newChange(tx, &change); err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Model(&or).Update("left_tier", relTier).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		if or.RightTier != relTier {
			change := Change{
				UserID:         userID,
				CompanyID:      or.RightID,
				RelationshipID: &or.ID,
				Type:           RelationshipUpdated,
				Key:            "tier",
				PreviousValue:  tiers[or.RightTier].Name,
				NewValue:       tiers[relTier].Name,
			}
			if err := newChange(tx, &change); err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Model(&or).Update("right_tier", relTier).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if or.Notes != notes {
		changeLeft := Change{
			UserID:         userID,
			CompanyID:      r.LeftID,
			RelationshipID: &or.ID,
			Type:           RelationshipUpdated,
			Key:            "notes",
			PreviousValue:  or.Notes,
			NewValue:       notes,
		}
		if err := newChange(tx, &changeLeft); err != nil {
			tx.Rollback()
			return err
		}
		changeRight := Change{
			UserID:         userID,
			CompanyID:      r.RightID,
			RelationshipID: &or.ID,
			Type:           RelationshipUpdated,
			Key:            "notes",
			PreviousValue:  or.Notes,
			NewValue:       notes,
		}
		if err := newChange(tx, &changeRight); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&or).Update("notes", notes).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

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
	if r.LeftID == r.RightID {
		return fmt.Errorf("a company can not be its own partner")
	}
	return nil
}

func (r *Relationship) EagerLoad(side int) error {
	switch side {
	case 1:
		c, err := GetCompany(r.LeftID, false)
		if err != nil {
			return err
		}
		r.LeftCompany = c
	case 2:
		c, err := GetCompany(r.RightID, false)
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

func (r *Relationship) Delete(userID uuid.UUID) error {
	if r.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}

	tx := session.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	changeLeft := Change{
		UserID:         userID,
		CompanyID:      r.LeftID,
		RelationshipID: &r.ID,
		Type:           RelationshipDeleted,
	}
	if err := newChange(tx, &changeLeft); err != nil {
		tx.Rollback()
		return err
	}

	changeRight := Change{
		UserID:         userID,
		CompanyID:      r.RightID,
		RelationshipID: &r.ID,
		Type:           RelationshipDeleted,
	}
	if err := newChange(tx, &changeRight); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(r).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func NewRelationship(r *Relationship, userID uuid.UUID) error {
	if r.ID == uuid.Nil {
		var err error
		if r.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if r.Conflicts() {
		return fmt.Errorf("the item already exists")
	}

	tx := session.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&r).Error; err != nil {
		tx.Rollback()
		return err
	}

	changeLeft := Change{
		UserID:         userID,
		CompanyID:      r.LeftID,
		RelationshipID: &r.ID,
		Type:           RelationshipCreated,
	}
	if err := newChange(tx, &changeLeft); err != nil {
		tx.Rollback()
		return err
	}

	changeRight := Change{
		UserID:         userID,
		CompanyID:      r.RightID,
		RelationshipID: &r.ID,
		Type:           RelationshipCreated,
	}
	if err := newChange(tx, &changeRight); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetRelationshipWithDeleted(id uuid.UUID) (*Relationship, error) {
	var item Relationship
	err := session.Unscoped().Where("id = ?", id).First(&item).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
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
