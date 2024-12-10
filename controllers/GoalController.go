package controllers

import (
	"fmt"
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

func (gc *GoalController) Index(c *gin.Context) {
	var goals []models.Goal

	// Debug: Start fetching goals
	fmt.Println("Fetching goals from the database...")
	err := gc.DB.Find(&goals).Error
	if err != nil {
		fmt.Printf("Error fetching goals: %v\n", err) // Debug: Print error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch goals"})
		return
	}

	// Debug: Log the fetched goals
	fmt.Printf("Goals fetched: %+v\n", goals)

	if len(goals) == 0 {
		fmt.Println("No goals found in the database.") // Debug: Empty database message
	}

	// Render the template
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Goal",
		"goals": goals,
	})
}

func (gc *GoalController) NewGoalSingle(c *gin.Context) {
	fmt.Println("Rendering the new goal form...") // Debug: Log rendering action

	var categories []models.Categories
	if err := gc.DB.Find(&categories).Error; err != nil {
		fmt.Printf("Error fetching categories: %v\n", err) // Debug: Print error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Debug: Log the fetched categories
	fmt.Printf("Categories fetched: %+v\n", categories)

	c.HTML(http.StatusOK, "goal.single.html", gin.H{
		"title":      "Add New Goal",
		"categories": categories,
	})
}

func (gc *GoalController) AddNewGoal(c *gin.Context) {
	fmt.Println("AddNewGoal endpoint hit...") // Debug: Log endpoint hit

	var goal models.Goal
	if err := c.ShouldBind(&goal); err != nil {
		fmt.Printf("Error binding input data: %v\n", err) // Debug: Print binding error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Debug: Log the bound goal object
	fmt.Printf("Goal object after binding: %+v\n", goal)

	var category models.Categories
	if err := gc.DB.Where("title = ?", goal.Category).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Category '%s' not found, creating new...\n", goal.Category) // Debug: Log new category creation

			category = models.Categories{Title: goal.Category}
			if err := gc.DB.Create(&category).Error; err != nil {
				fmt.Printf("Error creating category: %v\n", err) // Debug: Log category creation error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
				return
			}
		} else {
			fmt.Printf("Error fetching category: %v\n", err) // Debug: Log fetch error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
			return
		}
	}

	if err := gc.DB.Create(&goal).Error; err != nil {
		fmt.Printf("Error creating goal: %v\n", err) // Debug: Log goal creation error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	fmt.Printf("Goal successfully created: %+v\n", goal) // Debug: Log success
	c.Redirect(http.StatusSeeOther, "/")
}
