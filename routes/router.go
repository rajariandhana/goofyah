package routes

import (
	// "fmt"

	"fmt"
	"goofyah/controllers"
	"goofyah/middleware"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, store sessions.Store) *gin.Engine {
	router := gin.Default()
	router.Use(sessions.Sessions("this_session", store))

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
	unauthRoutes.Use(middleware.UnauthMiddleware())
	{
		unauthRoutes.GET("/login", authController.LoginCreate)
		unauthRoutes.POST("/login", authController.LoginStore)

		unauthRoutes.GET("/register", authController.RegisterCreate)
		unauthRoutes.POST("/register", authController.RegisterStore)
	}

	var list = []string{"anies", "prabowo", "ganjar"}
	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		authRoutes.POST("/logout", authController.LogoutStore)

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
			goalRoutes.GET("/addNewGoal", goalController.NewGoalSingle)
			goalRoutes.POST("/addNewGoal", goalController.AddGoal)
			// goalRoutes.GET("/:id", goalController.Show)
		}
		// }

		return router
	}
}
