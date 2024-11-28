package controllers

import (
	"net/http"
	// "goofyah/models"
	"goofyah/services"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	users := services.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

// Add CreateUser, GetUser, UpdateUser, DeleteUser functions
