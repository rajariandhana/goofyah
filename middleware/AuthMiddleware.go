package middleware

import (
	"fmt"
	"goofyah/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	// log.Print("auth")
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Auth")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid || claims["ttl"].(float64) < float64(time.Now().Unix()) {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		user, err := models.GetUserByID(uint(claims["userID"].(float64)))
		if err != nil {
			// log.Println("auth user not found", err)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		// log.Println("auth user found", user)

		c.Set("user", user)

		c.Next()
	}
}

func UnauthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Auth")
		if err != nil {
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		if claims["ttl"].(float64) < float64(time.Now().Unix()) {
			c.Next()
			return
		}

		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
