package routes

import (
	// "fmt"

	"fmt"
	"goofyah/controllers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	store := memstore.NewStore([]byte("secret"))
	router := gin.Default()
	router.Use(sessions.Sessions("temp", store))

	// router.LoadHTMLGlob("views/*.html")
	// router.LoadHTMLGlob("views/user/*.html")
	// router.LoadHTMLFiles("views/header.html", "views/footerr.html", "views/index.html")
	var htmlFiles []string
	err := filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file is an HTML file
		if filepath.Ext(path) == ".html" {
			htmlFiles = append(htmlFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error loading templates:", err)
		return nil
	}
	router.Static("/public", "./public")
	router.LoadHTMLFiles(htmlFiles...)
	authController := controllers.NewAuthController(db)
	goalController := controllers.NewGoalController(db)

	unauthRoutes := router.Group("/")
	// unauthRoutes.Use(middleware.UnauthMiddleware())
	// {
	unauthRoutes.GET("/login", authController.LoginCreate)
	unauthRoutes.POST("/login", authController.Login2)

	unauthRoutes.GET("/register", authController.RegisterCreate)
	unauthRoutes.POST("/register", authController.Register)
	// }

	var list = []string{"anies", "prabowo", "ganjar"}
	authRoutes := router.Group("/")
	// authRoutes.Use(middleware.AuthMiddleware3())
	// {
	authRoutes.POST("/logout", authController.Logout)

	authRoutes.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
			"list":  list,
		})
	})
	// contoh
	accountRoutes := authRoutes.Group("/account")
	{
		accountRoutes.GET("/", authController.Show)
	}
	goalRoutes := authRoutes.Group("/goals")
	{
		goalRoutes.GET("/", goalController.Index)
		// goalRoutes.GET("/:id", goalController.Show)
	}
	// }

	return router
}
