package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var session *gorm.DB

func Init(dbSession *gorm.DB) error {
	session = dbSession
	session.AutoMigrate(&User{})
	session.AutoMigrate(&Permission{})
	ok, err := rootUserCheck()
	if err != nil {
		return err
	}
	if !ok {
		if err := rootUserCreate(); err != nil {
			return err
		}
	}
	session.AutoMigrate(&Company{})
	session.AutoMigrate(&Relationship{})
	if err := UnmarshalCompanyTypes(); err != nil {
		return err
	}
	if err := UnmarshalSectors(); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("tiers", &tiers); err != nil {
		return err
	}
	return nil
}

func rootUserCheck() (bool, error) {
	users, err := ListUsers(nil)
	if err != nil {
		return false, err
	}
	if len(users) != 0 {
		return true, nil
	}
	return false, nil
}

func rootUserCreate() error {
	user := User{
		Username: viper.GetString("rootUser.username"),
		Password: viper.GetString("rootUser.password"),
		Email:    viper.GetString("rootUser.email"),
	}
	if err := NewUser(&user); err != nil {
		return err
	}
	pAdmin := Permission{
		UserID:     user.ID,
		ResourceID: WildCardResource,
		RoleID:     uuid.Must(uuid.Parse("23b0f96f-bd32-4421-96fa-e3bad618740c")),
	}
	return NewPermission(&pAdmin)
}
