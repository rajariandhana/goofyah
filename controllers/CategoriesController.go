package controllers

import (
	"goofyah/models"
	//"log"
	"net/http"
	"strconv"

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
	//for _, category := range categories {
	// log.Printf("Category: %+v\n", category)
	//log.Println("Category:", category.ID, category.Title, "user:", category.User.ID, category.User.Name)
	//}
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
	// categoryID := c.Param("ID")
	//log.Println("tes")
	categoryIDStr := c.Param("ID")
	//log.Println(categoryIDStr)

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	//log.Println(categoryID)
	category, err := models.GetCategoryByID(uint(categoryID))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Category not found"})
		return
	}

	goals := models.GetGoalsOfCategory(uint(categoryID))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch goals for category"})
	// 	return
	// }
	// log.Println(category.Goals)

	c.HTML(http.StatusOK, "single.category.html", gin.H{
		"title":    category.Title + " Goals",
		"category": category,
		"goals":    goals,
	})
}

func (gc *CategoriesController) DeleteCategories(c *gin.Context) {
	CategoryID := c.Param("ID")
	id, err := strconv.ParseUint(CategoryID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}
	var category models.Categories
	if err := gc.DB.Where("id = ?", uint(id)).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	if err := gc.DB.Where("categories_id = ?", uint(id)).Delete(&models.Goal{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goals"})
		return
	}
	if err := gc.DB.Delete(&models.Categories{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/categories/listcategories")
}
