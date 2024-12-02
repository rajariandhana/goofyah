package controllers

import (
	"goofyah/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf(user.Name + "|" + user.Email + "|" + user.Password)
	// check if email taken
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	user.Password = string(hash)
	if err := ac.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

type LoginInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginInput LoginInput

	if err := c.ShouldBindWith(&loginInput, binding.Form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// log.Printf("Email: %s, Password: %s\n", loginInput.Email, loginInput.Password)

	user, err := ac.GetUserByEmail(loginInput.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func (ac *AuthController) Login2(c *gin.Context) {
	var loginInput LoginInput
	var user models.User
	if err := ac.DB.Where("email = ?", loginInput.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (ac *AuthController) RegisterCreate(c *gin.Context) {
	users := ac.GetAllUser()
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
		"users": users,
	})
}

func (ac *AuthController) LoginCreate(c *gin.Context) {
	users := ac.GetAllUser()
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
		"users": users,
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) Logou2(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

func (ac *AuthController) Index(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("loggedIn") == nil || session.Get("loggedIn") != true {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (ac *AuthController) CreateUser(u *models.User) error {
	return ac.DB.Create(u).Error
}

func (ac *AuthController) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ac *AuthController) GetAllUser() []models.User {
	var users []models.User
	if err := ac.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return []models.User{}
	}
	log.Println("success get all")
	return users
}

func (ac *AuthController) Show(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}
	userObj, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User data type mismatch"})
		return
	}
	c.HTML(http.StatusOK, "user.show.html", gin.H{
		"title": "Login",
		"user":  userObj,
	})
}
