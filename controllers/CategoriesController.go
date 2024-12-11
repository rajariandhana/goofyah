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
		// log.Printf("Category: %+v\n", category)
		log.Println("Category:", category.ID, category.Title, "user:", category.User.ID, category.User.Name)
	}
}

func (uc *CategoriesController) Index(c *gin.Context) {
	models.ShowAllUser()
	uc.LogAllCategory()

	value, _ := c.Get("user")
	user := value.(*models.User)
	categories := models.GetCategoriesOfUser(*user)
	c.HTML(http.StatusOK, "categories.index.html", gin.H{
		"title":      "List of Categories",
		"categories": categories,
	})
}

type CategoryForm struct {
	Title string `form:"title"`
}

func (uc *CategoriesController) CreateCategory(c *gin.Context) {
	var form CategoryForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	value, _ := c.Get("user")
	user := value.(*models.User)
	// log.Println("uid", user.ID)
	var category models.Categories
	category.Title = form.Title
	category.UserID = user.ID
	category.User = *user

	if err := models.StoreCategories(category); err != nil {
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
