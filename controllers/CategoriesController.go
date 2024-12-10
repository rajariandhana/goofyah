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

func (uc *CategoriesController) CreateCategory(c *gin.Context) {
	var category models.Categories
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := uc.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save category"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/categories/listcategories")
}

func (uc *CategoriesController) ShowCategoryGoals(c *gin.Context) {

	categoryTitle := c.Param("category")

	var category models.Categories
	if err := uc.DB.Where("title = ?", categoryTitle).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var goals []models.Goal
	if err := uc.DB.Where("category = ?", categoryTitle).Find(&goals).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goals for this category not found"})
		return
	}

	c.HTML(http.StatusOK, "single.category.html", gin.H{
		"title":    category.Title + " Goals",
		"category": category,
		"goals":    goals,
	})
}
