package models

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var session *gorm.DB

func Init(dbSession *gorm.DB) error {
	session = dbSession

	// Authentication bits
	session.AutoMigrate(&User{})
	session.AutoMigrate(&Permission{})
	session.AutoMigrate(&Session{})
	if err := rootUserInit(); err != nil {
		return err
	}

	// Application bits
	session.AutoMigrate(&Company{})
	session.AutoMigrate(&CompanyVertical{})
	session.AutoMigrate(&Relationship{})
	session.AutoMigrate(&Change{})
	if err := UnmarshalCompanyTypes(); err != nil {
		return err
	}
	if err := UnmarshalVerticals(); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("tiers", &tiers); err != nil {
		return err
	}
	return nil
}
