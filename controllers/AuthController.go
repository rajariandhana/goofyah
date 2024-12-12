package controllers

import (
	"goofyah/models"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type FormInput struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Address  string `form:"address"`
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
		user, err := models.GetUserByEmail(input.Email)
		if err == nil && user != nil {
			err_email = "Email already taken"
		}
	}

	var hash []byte
	if input.Password == "" {
		err_password = "Invalid password"
	} else {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			err_password = "Something wrong when storing password"
		} else {
			hash = hashed
		}
	}

	if err_name == "" && err_email == "" && err_password == "" {
		var user models.User
		user.Name = input.Name
		user.Email = input.Email
		user.Password = string(hash)
		user.Address = input.Address
		if err := models.StoreUser(user); err == nil {
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
	// models.ShowAllUser()
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
	err_fail := ""

	err_email = ac.EmailValid(input.Email)
	// log.Print(err_email)
	var user models.User
	if err_email == "" {
		existingUser, err := models.GetUserByEmail(input.Email)
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
		// log.Println("hashed pass ", []byte(user.Password))
		// hash, x := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		// if x != nil && hash == nil {
		// }
		// log.Println("hashed input ", hash)
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			err_password = "Wrong password"
			log.Print(err_password)
		}
	}

	if err_email == "" && err_password == "" {
		tokenString, err := createAndSignJWT(&user)
		if err != nil {
			err_fail = "Failed to log in"
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"title":    "Login",
				"err_fail": err_fail,
			})
			return
		}
		// log.Println("tokenString ", tokenString)
		// log.Println("error nil")
		setCookie(c, tokenString)
		// log.Println("cookie set redirecting to /")
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusBadRequest, "login.html", gin.H{
		"title":        "Login",
		"err_email":    err_email,
		"err_password": err_password,
		"err_fail":     err_fail,
	})

}

func createAndSignJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"ttl":    time.Now().Add(time.Hour * 24 * 100).Unix(),
	})
	// log.Println("secret ", []byte(os.Getenv("SECRET")))
	secret := strings.TrimSpace(os.Getenv("SECRET"))
	return token.SignedString([]byte(secret))
}

func setCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*100, "", "", false, true)
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
	c.SetCookie("Auth", "deleted", 0, "", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (ac *AuthController) Show(c *gin.Context) {
	user, _ := c.Get("user")
	// log.Println("showing user ", user)
	models.ShowAllUser()
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

	value, _ := c.Get("user")
	user := value.(*models.User)

	if input.Name == "" {
		err_name = "Invalid name"
	} else {
		user.Name = input.Name
	}

	err_email = ac.EmailValid(input.Email)
	if err_email == "" && input.Email != user.Email {
		existingUser, err := models.GetUserByEmail(input.Email)
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
		if models.SaveUser(*user) == nil {
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
