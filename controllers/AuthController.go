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

type FormInput struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (ac *AuthController) RegisterCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

func (ac *AuthController) RegisterStore(c *gin.Context) {
	var input FormInput
	if err := c.ShouldBind(&input); err != nil {
		log.Println(err.Error())
		return
	}

	err_name := ""
	err_email := ""
	err_password := ""
	err_storing := ""
	var user models.User

	if input.Name == "" {
		err_name = "Invalid name"
	} else {
		user.Name = input.Name
	}

	mes := ac.EmailValid(input.Email, "")
	if mes != "" {
		err_email = mes
	} else {
		user.Email = input.Email
	}

	if input.Password == "" {
		err_password = "Invalid password"
	} else if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {

		}
		user.Password = string(hash)
	}

	if err_name == "" || err_email == "" || err_password == "" {
		if err := ac.DB.Create(&user).Error; err == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		err_storing = "Failed to create user"
	}

	c.HTML(http.StatusBadRequest, "user.show.html", gin.H{
		"title":        "user.show",
		"user":         user,
		"err_storing":  err_storing,
		"err_email":    err_email,
		"err_name":     err_name,
		"err_password": err_password,
	})
}

func (ac *AuthController) LoginCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func (ac *AuthController) LoginStore(c *gin.Context) {
	var input FormInput

	if err := c.ShouldBind(&input); err != nil {
		log.Println(err.Error())
		return
	}

	err_email := ""
	err_password := ""

	var user models.User

	if input.Password==""{
		err_password = "Password cannot be empty"
	}
	
	mes := ac.EmailValid(input.Email, "")
	log.Println("email is " + input.Email)
	if mes != "" {
		err_email = mes
	} else {
		log.Println("email valid")
		user, err := ac.GetUserByEmail(input.Email)

		if err == nil && user != nil {
			user.Email = input.Email
		} else {
			err_email = "Failed to get email"
		}
	}

	 else if{

	}

	if err_email == "" && err_password == "" {
		session := sessions.Default(c)
		session.Set("userID", user.ID)
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusBadRequest, "login.html", gin.H{
		"title":        "Login",
		"user":         user,
		"err_email":    err_email,
		"err_password": err_password,
	})

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
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := ac.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "user.show.html", gin.H{
		"title": "Login",
		"user":  user,
	})
}

func (ac *AuthController) Update(c *gin.Context) {
	var updateInput FormInput
	if err := c.ShouldBind(&updateInput); err != nil {
		log.Println(err.Error())
		return
	}

	err_name := ""
	err_email := ""
	err_updating := ""

	session := sessions.Default(c)
	userID := session.Get("userID")
	user, err := ac.GetUserByID(userID.(uint))
	if err != nil {
		err_updating = "Failed to update user"
	}

	if updateInput.Name == "" {
		err_name = "Invalid name"
	} else {
		user.Name = updateInput.Name
	}

	mes := ac.EmailValid(updateInput.Email, user.Email)
	if mes != "" {
		err_email = mes
	} else {
		user.Email = updateInput.Email
	}

	if updateInput.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(updateInput.Password), bcrypt.DefaultCost)
		if err != nil {

		}
		user.Password = string(hash)
	}

	if err_updating == "" && err_email == "" && err_name == "" {
		if err := ac.DB.Save(&user).Error; err == nil {
			success := "Information successfully updated"
			c.HTML(http.StatusFound, "user.show.html", gin.H{
				"title":   "user.show",
				"user":    user,
				"success": success,
			})
			return
		}
		err_updating = "Failed to update user"
	}

	c.HTML(http.StatusBadRequest, "user.show.html", gin.H{
		"title":        "user.show",
		"user":         user,
		"err_updating": err_updating,
		"err_name":     err_name,
		"err_email":    err_email,
	})
}

func (ac *AuthController) EmailValid(email string, userEmail string) string {
	log.Println("emailvalid " + email)
	re := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(re, email)
	msg := ""
	if !match {
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
