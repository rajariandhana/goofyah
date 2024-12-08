package models

import (
	"gorm.io/gorm"
)

type Categories struct {
	gorm.Model
	// ID       uint   `gorm:"primaryKey" json:"id"`
	Title string `gorm:"size:64" json:"title" form:"title" binding:"required"`
}

// jangan lupa panggil MigrateGoal di connection.go
/*func MigrateGoal(db *gorm.DB) {
	db.AutoMigrate(&Categories{})
}
*/
