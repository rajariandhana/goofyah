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
	categoriesController := controllers.NewCategoriesController(db)

	unauthRoutes := router.Group("/")
	unauthRoutes.Use(middleware.UnauthMiddleware())
	{
		unauthRoutes.GET("/login", authController.LoginCreate)
		unauthRoutes.POST("/login", authController.LoginStore)

		unauthRoutes.GET("/register", authController.RegisterCreate)
		unauthRoutes.POST("/register", authController.RegisterStore)
	}

	authRoutes := router.Group("/")

	authRoutes.Use(middleware.AuthMiddleware())
	{
		authRoutes.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Home",
			})
		})

		// contoh
		accountRoutes := authRoutes.Group("/account")
		{
			accountRoutes.GET("", authController.Show)
			accountRoutes.POST("/update", authController.Update)
			accountRoutes.POST("/logout", authController.LogoutStore)
		}
		goalRoutes := authRoutes.Group("goals")
		{
			goalRoutes.GET("", goalController.Index)
			goalRoutes.GET("/addNewGoal", goalController.NewGoalSingle)
			goalRoutes.POST("/addNewGoal", goalController.AddGoal)
			// goalRoutes.GET("/:id", goalController.Show)
		}
		categoriesRoutes := authRoutes.Group("/categories")
		{
			categoriesRoutes.GET("/listcategories", categoriesController.Index)           // GET request to fetch and display categories at /categories/listcategories
			categoriesRoutes.POST("/listcategories", categoriesController.CreateCategory) // POST request to create a new category at /categories/listcategories

		}
	}
	// }

	return router
}
