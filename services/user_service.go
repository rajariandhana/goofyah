package services

import (
	"goofyah/database"
	"goofyah/models"
)

func GetAllUsers() []models.User {
	var users []models.User
	database.DB.Find(&users)
	return users
}
