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

	if input.Name == "" {
		err_name = "Invalid name"
	}

	err_email = ac.EmailValid(input.Email)
	if err_email == "" {
		user, err := ac.GetUserByEmail(input.Email)
		if err == nil && user != nil {
			err_email = "Email already taken"
		}
	}

	var hash []byte
	if input.Password == "" {
		err_password = "Invalid password"
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil && hash == nil {
			err_password = "Something wrong when storing password"
		}
	}

	if err_name == "" && err_email == "" && err_password == "" {
		var user models.User
		user.Name = input.Name
		user.Email = input.Email
		user.Password = string(hash)
		if err := ac.DB.Create(&user).Error; err == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		err_storing = "Failed to create user"
	}

	c.HTML(http.StatusBadRequest, "register.html", gin.H{
		"title":        "Register",
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

	err_email = ac.EmailValid(input.Email)
	log.Print(err_email)
	var user models.User
	if err_email == "" {
		existingUser, err := ac.GetUserByEmail(input.Email)
		if existingUser == nil && err != nil {
			err_email = "Email not found"
			log.Print(err_email)
		} else {
			user = *existingUser
		}
	}

	if input.Password == "" {
		err_password = "Password cannot be empty"
		log.Print(err_password)
	}

	if err_email == "" && err_password == "" {
		log.Println("user found " + user.Email)
		// Hashed(user.Password, input.Password)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			err_password = "Wrong password"
			log.Print(err_password)
		}
	}

	if err_email == "" && err_password == "" {
		session := sessions.Default(c)
		log.Println("user.ID ", user.ID)
		session.Set("userID", user.ID)
		// log.Println("session.Get(userID) ", session.Get("userID"))
		session.Save()
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusBadRequest, "login.html", gin.H{
		"title":        "Login",
		"err_email":    err_email,
		"err_password": err_password,
	})

}

// func Hashed(pass1 string, pass2 string) {
// 	hp1, err1 := bcrypt.GenerateFromPassword([]byte(pass1), bcrypt.DefaultCost)
// 	hp2, err2 := bcrypt.GenerateFromPassword([]byte(pass2), bcrypt.DefaultCost)
// 	if err1 == nil && err2 == nil {
// 		log.Println(hp1)
// 		log.Println(hp2)
// 	}
// }

func (ac *AuthController) LogoutStore(c *gin.Context) {
	log.Println("LogoutStore")
	session := sessions.Default(c)
	// log.Println("userid ", session.Get("userID"))
	session.Clear()
	// log.Println("userid ", session.Get("userID"))
	session.Save()
	// log.Println("userid ", session.Get("userID"))
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
		"title": "Account",
		"user":  user,
	})
}

func (ac *AuthController) Update(c *gin.Context) {
	var input FormInput
	if err := c.ShouldBind(&input); err != nil {
		log.Println(err.Error())
		return
	}

	err_name := ""
	err_email := ""
	err_password := ""
	err_updating := ""

	session := sessions.Default(c)
	userID := session.Get("userID")
	user, err := ac.GetUserByID(userID.(uint))
	if err != nil {
		err_updating = "Failed to update user"
	}

	if input.Name == "" {
		err_name = "Invalid name"
	} else {
		user.Name = input.Name
	}

	err_email = ac.EmailValid(input.Email)
	if err_email == "" && input.Email != user.Email {
		existingUser, err := ac.GetUserByEmail(input.Email)
		if err == nil && existingUser != nil {
			err_email = "Email already taken"
		} else {
			user.Email = input.Email
		}
	}

	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			err_password = "Something wrong when storing password"
		}
		user.Password = string(hash)
	}

	if err_updating == "" && err_email == "" && err_name == "" && err_password == "" {
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
		"err_password": err_email,
	})
}

func (ac *AuthController) EmailValid(email string) string {
	// log.Println("emailvalid " + email)
	re := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(re, email)
	if !match {
		return "Invalid email"
	}
	return ""
}

/*
if empty
name cannot be empty
email is not valid
email is taken

minimum of 8 characters
*/
