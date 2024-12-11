package models

import (
	//"time" //ini kl mau include timestamp (could be waktu / tanggal)

	"gorm.io/gorm"
)

type Goal struct {
	gorm.Model
	Title string `gorm:"size:64" json:"title" form:"title" binding:"required"` // Changed GoalTitle to Title
	// Category    string `gorm:"size:255" json:"category" form:"category" binding:"required"`
	Description string `gorm:"size:512" json:"description" form:"description"`
	// ID       uint   `gorm:"primaryKey" json:"id"`
	//StartAt     time.Time `json:"startat" form:"startat" binding:"required"`
	//EndAt       time.Time `json:"endat" form:"endat" binding:"required"`
	// UserID       uint       `gorm:""`
	// User         User       `gorm:"constraint:OnDelete:CASCADE;"`
	CategoriesID uint       `gorm:""`
	Categories   Categories `gorm:"constraint:OnDelete:CASCADE;"`
}

func GetGoalsOfUser(user User) []Goal {
	var goals []Goal
	DB.Preload("Categories.Goals").First(&user, user.ID)
	for _, category := range user.Categories {
		for _, goal := range category.Goals {
			// Associate the category with each goal for display
			goal.Categories = category
			goals = append(goals, goal)
		}
	}
	return goals
}

func StoreGoal(goal Goal) error {
	return DB.Create(&goal).Error
}
