package models

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// ID         uint         `gorm:"primaryKey" json:"id"`
	Name       string       `gorm:"size:64" json:"name" form:"name" binding:"required"`
	Email      string       `gorm:"size:64,index" json:"email" form:"email" binding:"required,email"`
	Password   string       `gorm:"size:255" json:"password" form:"password" binding:"required"`
	Address    string       `gorm:"" form:"address"`
	Categories []Categories `gorm:"foreignKey:UserID"`
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

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		// log.Println("email not exist")
		return nil, err
	}
	return &user, nil
}

func GetAllUser() []User {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return []User{}
	}
	// log.Println("success get all")
	return users
}

func ShowAllUser() {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return
	}
	for _, user := range users {
		// log.Println(user)
		log.Println("user", user.ID, user.Name, user.Email, user.Address)
	}
}

func StoreUser(user User) error {
	return DB.Create(&user).Error
}

func SaveUser(user User) error {
	return DB.Save(&user).Error
}
