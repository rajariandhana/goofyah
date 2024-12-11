package controllers

import (
	"fmt"
	"goofyah/models"
	"log"
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
	log.Println("indexgoal")
	gc.ShowAllGoal()

	value, _ := c.Get("user")
	user := value.(*models.User)

	var goals []models.Goal
	gc.DB.Preload("Categories.Goals").First(&user, user.ID)
	for _, category := range user.Categories {
		goals = append(goals, category.Goals...)
	}

	log.Println("Goals fetched:\n", goals)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Goal",
		"goals": goals,
	})
}

func (gc *GoalController) ShowAllGoal() {
	var goals []models.Goal
	if err := gc.DB.Find(&goals).Error; err != nil {
		return
	}
	for _, goal := range goals {
		log.Println(goal)
	}
}

func (gc *GoalController) NewGoalSingle(c *gin.Context) {
	fmt.Println("Rendering the new goal form...")
	value, _ := c.Get("user")
	user := value.(*models.User)
	gc.DB.Preload("Categories").First(&user, user.ID)

	c.HTML(http.StatusOK, "goal.single.html", gin.H{
		"title":      "Add New Goal",
		"categories": user.Categories,
	})
}

type GoalForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	CategoryID  uint   `form:"category"`
	// CategoryName string `form:"categoryName"`
}

func (gc *GoalController) AddNewGoal(c *gin.Context) {
	fmt.Println("AddNewGoal endpoint hit...") // Debug: Log endpoint hit
	var form GoalForm
	var goal models.Goal
	if err := c.ShouldBind(&form); err != nil {
		fmt.Printf("Error binding input data: %v\n", err) // Debug: Print binding error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var category models.Categories
	if err := gc.DB.Where("ID = ?", form.CategoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
		// log.Println("emm not found categoryid")
		// value, _ := c.Get("user")
		// user := value.(*models.User)
		// category.Title = form.CategoryName
		// category.UserID = user.ID
		// category.User = *user
		// if err := gc.DB.Create(&category).Error; err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save category"})
		// 	return
		// }
	}
	goal.Title = form.Title
	goal.Description = form.Description
	goal.CategoriesID = category.ID
	goal.Categories = category

	if err := gc.DB.Create(&goal).Error; err != nil {
		fmt.Printf("Error creating goal: %v\n", err) // Debug: Log goal creation error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	fmt.Printf("Goal successfully created: %+v\n", goal) // Debug: Log success
	c.Redirect(http.StatusSeeOther, "/")
}
