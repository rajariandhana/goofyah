package controllers

import (
	"goofyah/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoriesController struct {
	DB *gorm.DB
}

func NewCategoriesController(db *gorm.DB) *CategoriesController {
	return &CategoriesController{DB: db}
}

func (uc *CategoriesController) LogAllCategory() {
	var categories []models.Categories
	if err := uc.DB.Find(&categories).Error; err != nil {
		return
	}
	for _, category := range categories {
		log.Printf("Category: %+v\n", category)
	}
}

// Display all categories
func (uc *CategoriesController) Index(c *gin.Context) {
	uc.LogAllCategory()
	var categories []models.Categories
	if err := uc.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categories not found"})
		return
	}

	c.HTML(http.StatusOK, "categories.index.html", gin.H{
		"title":      "List of Categories",
		"categories": categories,
	})
}

// Handle category creation
func (uc *CategoriesController) CreateCategory(c *gin.Context) {
	var category models.Categories
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Save the category to the database
	if err := uc.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save category"})
		return
	}

	// Redirect to the categories page to display the updated list
	c.Redirect(http.StatusSeeOther, "/categories/listcategories")
}
