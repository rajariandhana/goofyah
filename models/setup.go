package models

import "gorm.io/gorm"

var DB *gorm.DB

func ConnectDB(db *gorm.DB) {
	DB = db
}
