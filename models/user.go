package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:64" json:"name" form:"name" binding:"required"`
	Email    string `gorm:"size:64,index" json:"email" form:"email" binding:"required,email"`
	Password string `gorm:"size:255" json:"password" form:"password" binding:"required"`
}

func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

func GetUserByID(id uint) (*User, error) {
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}
