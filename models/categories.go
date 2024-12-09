package models

import (
	"gorm.io/gorm"
)

type Categories struct {
	gorm.Model
	Title string `gorm:"size:64" json:"title" form:"title" binding:"required"`
}
