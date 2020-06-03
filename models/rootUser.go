package models

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func rootUserInit() error {
	ok, err := rootUserCheck()
	if err != nil {
		return err
	}
	if !ok {
		if err := rootUserCreate(); err != nil {
			return err
		}
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
		RealName: viper.GetString("rootUser.realname"),
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
