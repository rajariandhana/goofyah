package middleware

import (
	"fmt"
	"goofyah/models"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	log.Print("auth")
	return func(c *gin.Context) {
		log.Print("authx")
		tokenString, err := c.Cookie("Auth")
		log.Println("tokenString ", tokenString)
		if err != nil {
			log.Println("err1")
			log.Println(err)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("err2")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// log.Println("secret ", []byte(os.Getenv("SECRET")))
			secret := strings.TrimSpace(os.Getenv("SECRET"))
			return []byte(secret), nil
		})
		if err != nil {
			log.Println("err3")
			log.Println(err)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid || claims["ttl"].(float64) < float64(time.Now().Unix()) {
			log.Println("err4")
			log.Println("ok: ", ok)
			log.Println("token.Valid: ", token.Valid)
			log.Println("ttl: ", claims["ttl"].(float64) < float64(time.Now().Unix()))
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		user, err := models.GetUserByID(uint(claims["userID"].(float64)))
		if err != nil {
			log.Println("auth user not found", err)
			log.Println("err5")
			log.Println(err)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		log.Println("auth user found", user)

		c.Set("user", user)

		c.Next()
	}
}

func UnauthMiddleware() gin.HandlerFunc {
	log.Println("unauht")
	return func(c *gin.Context) {
		log.Print("unauthx")
		tokenString, err := c.Cookie("Auth")
		if err != nil {
			log.Println("err1")
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			log.Println("err2")
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			log.Println("err3")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			log.Println("err4")
			return
		}

		if claims["ttl"].(float64) < float64(time.Now().Unix()) {
			c.Next()
			log.Println("err5")
			return
		}

		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			log.Println("err6")
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		log.Println("err7")
		c.Next()
	}
}
