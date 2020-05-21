package models

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var session *gorm.DB

func Init(dbSession *gorm.DB) error {
	session = dbSession
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
