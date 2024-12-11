package models

import (
	"gorm.io/gorm"
)

type Categories struct {
	gorm.Model
	Title  string `gorm:"size:64" json:"title" form:"title" binding:"required"`
	UserID uint   `gorm:""`
	User   User   `gorm:"constraint:OnDelete:CASCADE;"`
	Goals  []Goal `gorm:"foreignKey:CategoriesID"`
}

func GetCategoriesOfUser(user User) []Categories {
	DB.Preload("Categories").First(&user, user.ID)
	return user.Categories
}

func StoreCategories(categories Categories) error {
	return DB.Create(&categories).Error
}

func GetCategoryByID(ID uint) (*Categories, error) {
	var category Categories
	if err := DB.First(&category, ID).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func GetGoalsOfCategory(ID uint) []Goal {
	var goals []Goal
	DB.Preload("Goals").First(&goals, ID)
	return goals
}
