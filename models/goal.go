package models

import "gorm.io/gorm"

type Goal struct {
}

// jangan lupa panggil MigrateGoal di connection.go
func MigrateGoal(db *gorm.DB) {
	db.AutoMigrate(&Goal{})
}
