package controllers

import (
	"goofyah/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GoalController struct {
	DB *gorm.DB
}

func NewGoalController(db *gorm.DB) *GoalController {
	return &GoalController{DB: db}
}

func (uc *GoalController) Index(c *gin.Context) {
	var goals []models.Goal
	if err := uc.DB.Find(&goals).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "user.index.html", gin.H{
		"title": "Goal",
		"goals": goals,
	})
}
