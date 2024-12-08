package controllers

import (
	"goofyah/models"
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

func (uc *CategoriesController) Index(c *gin.Context) {
	var categories []models.Categories
	if err := uc.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Categories not found"})
		return
	}

	c.HTML(http.StatusOK, "categories.index.html", gin.H{
		"title":      "categories",
		"categories": categories,
	})
}

func (gc *CategoriesController) CategoriesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "categories.index.html", gin.H{
		"title": "List of Categories",
	})
}
