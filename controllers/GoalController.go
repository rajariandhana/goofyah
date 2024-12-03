package controllers

import (
	"goofyah/models"
	"net/http"

	"fmt"

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
		c.JSON(http.StatusNotFound, gin.H{"error": "Goals not found"})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Goal",
		"goals": goals,
	})
}

func (gc *GoalController) NewGoalSingle(c *gin.Context) {
	fmt.Println("Rendering the new goal form") // Debugging log

	c.HTML(http.StatusOK, "goal.single.html", gin.H{
		"title": "Add New Goal",
	})
}

func (gc *GoalController) AddGoal(c *gin.Context) {
	var goal models.Goal

	if err := c.ShouldBind(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := gc.DB.Create(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Goal created successfully",
		"goal":    goal,
	})
}
