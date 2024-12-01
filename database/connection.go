package database

import (
	"goofyah/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	// var err error
	DB, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	return DB, nil
}
