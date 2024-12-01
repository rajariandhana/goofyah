package controllers

import (

	// "goofyah/models"

	"fmt"
	"goofyah/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) Index(c *gin.Context) {
	var users []models.User
	if err := uc.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// c.JSON(200, gin.H{"users": users})
	c.HTML(http.StatusOK, "user.index.html", gin.H{
		"title": "User",
		"users": users,
	})
}

func (uc *UserController) Show(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.HTML(http.StatusOK, "user.show.html", gin.H{
		"title": "User",
		"user":  user,
	})
}

func (uc *UserController) Create(c *gin.Context) {
	c.HTML(http.StatusOK, "user.create.html", gin.H{
		"title": "Create User",
	})
}

func (uc *UserController) Store(c *gin.Context) {
	var user models.User
	fmt.Printf(user.Name)
	fmt.Printf(user.Email)
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func (uc *UserController) Update(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind data"})
		return
	}
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

func (uc *UserController) Destroy(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := uc.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.Redirect(http.StatusFound, "/users")
}

// GetUser retrieves a user by ID

// Add CreateUser, GetUser, UpdateUser, DeleteUser functions

// user := models.User{Name: "John Doe", Email: "john@example.com"}
// result := uc.DB.Create(&user)
// if result.Error != nil {
// 	fmt.Println("Error creating user:", result.Error)
// 	return
// }
