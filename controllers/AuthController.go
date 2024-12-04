package controllers

import (
	"errors"
	"fmt"
	"goofyah/models"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) RegisterCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (ac *AuthController) RegisterStore(c *gin.Context) {
	var user models.User
	err_n := 0
	err_invalid_input := ""
	err_email_taken := ""
	err_email_fail := ""
	err_user_fail_create := ""

	if err := c.ShouldBind(&user); err != nil {
		err_invalid_input = "Invalid input data"
		err_n++
	}
	// log.Printf(user.Name + "|" + user.Email + "|" + user.Password)
	existingUser, err := ac.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		err_email_taken = "Email is already taken"
		err_n++
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		err_email_fail = "Failed to check email"
		err_n++
	}

	if err_n == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			err_user_fail_create = "Failed to create user"
			err_n++
		}
		user.Password = string(hash)
		if err = ac.CreateUser(&user); err != nil {
			err_user_fail_create = "Failed to create user"
			err_n++
		}
	}

	if err_n > 0 {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title":                "Register",
			"user":                 user,
			"err_invalid_input":    err_invalid_input,
			"err_email_taken":      err_email_taken,
			"err_email_fail":       err_email_fail,
			"err_user_fail_create": err_user_fail_create,
		})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) LoginCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

type LoginInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

func (ac *AuthController) LoginStore(c *gin.Context) {
	var loginInput LoginInput
	err_n := 0
	err_invalid_input := ""
	err_email_not_found := ""
	err_password_wrong := ""
	if err := c.ShouldBind(&loginInput); err != nil {
		err_invalid_input = "Invalid input data"
		err_n++
	}

	var user models.User
	if err := ac.DB.Where("email = ?", loginInput.Email).First(&user).Error; err != nil {
		err_email_not_found = "Email not found"
		err_n++
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password)); err != nil {
		err_password_wrong = "Wrong password"
		err_n++
	}

	if err_n > 0 {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title":               "Login",
			"user":                user,
			"err_invalid_input":   err_invalid_input,
			"err_email_not_found": err_email_not_found,
			"err_password_wrong":  err_password_wrong,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (ac *AuthController) LogoutStore(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) CreateUser(u *models.User) error {
	return ac.DB.Create(u).Error
}

func (ac *AuthController) GetAllUser() []models.User {
	var users []models.User
	if err := ac.DB.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return []models.User{}
	}
	// log.Println("success get all")
	return users
}

func (ac *AuthController) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ac.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ac *AuthController) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := ac.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (ac *AuthController) Show(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	log.Println("er1")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	log.Println("er2")
	user, err := ac.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	log.Println("er3")
	c.HTML(http.StatusOK, "user.show.html", gin.H{
		"title": "Login",
		"user":  user,
	})
}

type UpdateInput struct {
	Name     string `form:"name"`
	Email    string `form:"email" binding:"email"`
	Password string `form:"password"`
}

func (ac *AuthController) Update(c *gin.Context) {
	var updateInput UpdateInput
	if err := c.ShouldBind(&updateInput); err != nil {
		log.Println(err.Error())
		return
	}

	err_n := 0
	err_name := ""
	err_email := ""
	err_updating := ""

	session := sessions.Default(c)
	userID := session.Get("userID")
	user, err := ac.GetUserByID(userID.(uint))
	if err != nil {
		// c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		err_updating = "Failed to update user"
		err_n++
	}

	if updateInput.Name == "" {
		err_name = "Invalid name"
		err_n++
	} else {
		user.Name = updateInput.Name
	}

	mes := ac.EmailValid(updateInput.Email, user.Email)
	if mes != "" {
		err_email = mes
		err_n++
	} else {
		user.Email = updateInput.Email
	}

	if updateInput.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(updateInput.Password), bcrypt.DefaultCost)
		if err != nil {

		}
		user.Password = string(hash)
	}

	if err_n > 0 {
		c.HTML(http.StatusBadRequest, "user.show.html", gin.H{
			"title":        "user.show",
			"user":         user,
			"err_updating": err_updating,
			"err_name":     err_name,
			"err_email":    err_email,
		})
		return
	}

	if err := ac.DB.Save(&user).Error; err != nil {
		err_updating = "Failed to update user"
		c.HTML(http.StatusBadRequest, "user.show.html", gin.H{
			"title":        "user.show",
			"user":         user,
			"err_updating": err_updating,
		})
		return
	}

	success := "Information successfully updated"
	c.HTML(http.StatusFound, "user.show.html", gin.H{
		"title":   "user.show",
		"user":    user,
		"success": success,
	})
}

func (ac *AuthController) EmailValid(email string, userEmail string) string {
	re := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(re, email)
	msg := ""
	if match == false {
		msg = "Invalid email"
	} else if email != userEmail {
		existingUser, err := ac.GetUserByEmail(userEmail)
		if err == nil && existingUser != nil {
			msg = "Email is already taken"
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			msg = "Failed to check email"
		}
	}
	return msg
}

/*
if empty
name cannot be empty
email is not valid
email is taken

minimum of 8 characters
*/
