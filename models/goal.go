package models

import (
	"time" //ini kl mau include timestamp (could be waktu / tanggal)

	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	// ID       uint   `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:64" json:"title" form:"title" binding:"required"`
	StartAt     time.Time `json:"startat" form:"startat" binding:"required"`
	EndAt       time.Time `json:"endat" form:"endat" binding:"required"`
	Description string    `gorm:"size:512" json:"description" form:"description"`
}
