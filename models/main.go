package models

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var session *gorm.DB

var Tiers []Tier

func Init(dbSession *gorm.DB) error {
	session = dbSession
	session.AutoMigrate(&Company{})
	session.AutoMigrate(&CompanyType{})
	session.AutoMigrate(&Relationship{})
	session.AutoMigrate(&Sector{})
	err := viper.UnmarshalKey("tiers", &Tiers)
	if err != nil {
		return err
	}

	return nil
}
