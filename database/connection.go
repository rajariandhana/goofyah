package database

import (
	"goofyah/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() (*gorm.DB, error) {
	var err error
	DB, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect" + err.Error())
	}
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}
	log.Println("Database connected successfully")
	return DB, nil
}
