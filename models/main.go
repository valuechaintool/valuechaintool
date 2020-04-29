package models

import "github.com/jinzhu/gorm"

var session *gorm.DB

func Init(dbSession *gorm.DB) error {
	session = dbSession
	session.AutoMigrate(&Company{})
	session.AutoMigrate(&Relationship{})
	return nil
}
