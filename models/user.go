package models

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `form:"name" json:"name" binding:"required"`
	Email string `form:"email" json:"email" binding:"required,email"`
}
