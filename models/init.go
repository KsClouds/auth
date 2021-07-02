package models

import (
	"auth/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func StartMySQL() (err error) {
	DB, err = gorm.Open(config.GetConfigInf())
	DB.AutoMigrate(&User{}, &Superior{})
	DB.LogMode(true)
	return err
}
