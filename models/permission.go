package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var WildCardResource = uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000001"))

type Permission struct {
	ID         uuid.UUID  `json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time ``
	UserID     uuid.UUID  `json:"user_id"`
	ResourceID uuid.UUID  `json:"resource_id"`
	RoleID     uuid.UUID  `json:"role_id"`
}

func (p *Permission) BeforeSave() error {
	if err := p.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *Permission) Conflicts() bool {
	return false
}

func (p *Permission) Validate() error {
	return nil
}

func (p *Permission) Delete() error {
	if p.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}
	return session.Delete(p).Error
}

func NewPermission(p *Permission) error {
	if p.ID == uuid.Nil {
		var err error
		if p.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if p.Conflicts() {
		return fmt.Errorf("the item already exists")
	}
	if err := session.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func GetPermission(id uuid.UUID) (*Permission, error) {
	var item Permission
	err := session.Where("id = ?", id).First(&item).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func ListPermissionsByUser(userID uuid.UUID) ([]Permission, error) {
	var items []Permission
	if err := session.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func ListCapabilitiesByUser(userID uuid.UUID) (map[uuid.UUID][]string, error) {
	permissions, err := ListPermissionsByUser(userID)
	if err != nil {
		return nil, err
	}
	caps := make(map[uuid.UUID][]string)
	for _, permission := range permissions {
		role, err := GetRole(permission.RoleID)
		if err != nil {
			return nil, err
		}
		if _, ok := caps[permission.ResourceID]; !ok {
			caps[permission.ResourceID] = role.Capabilities
		} else {
			caps[permission.ResourceID] = append(caps[permission.ResourceID], role.Capabilities...)
		}
	}
	return caps, nil
}
