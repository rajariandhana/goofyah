package middleware

import (
	"fmt"
	"goofyah/database"
	"goofyah/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			var user models.User
			database.DB.First(&user, claims["sub"])
			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Set("user", user)
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
func AuthMiddleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func UnauthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Authorization")
		if err == nil && token != "" {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthMiddleware3() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("userID")
		if sessionID == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		userID := sessionID.(uint)
		var user models.User
		database.DB.Where("id = ?", userID).First(&user)
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
