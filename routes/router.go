package routes

import (
	"fmt"
	"goofyah/controllers"
	"goofyah/middleware"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

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
	authController := controllers.NewAuthController()
	goalController := controllers.NewGoalController(db)
	categoriesController := controllers.NewCategoriesController(db)

	router.GET("/login", middleware.UnauthMiddleware(), authController.LoginCreate)
	router.POST("/login", middleware.UnauthMiddleware(), authController.LoginStore)

	router.GET("/register", middleware.UnauthMiddleware(), authController.RegisterCreate)
	router.POST("/register", middleware.UnauthMiddleware(), authController.RegisterStore)

	router.Use(middleware.AuthMiddleware())
	// all routes below will use AuthMiddleware
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
		})
	})

	userRoutes := router.Group("/user")
	{
		userRoutes.GET("", authController.Show)
		userRoutes.POST("/update", authController.Update)
		userRoutes.POST("/logout", authController.LogoutStore)
	}
	goalRoutes := router.Group("/goals")
	{
		goalRoutes.GET("", goalController.Index)
		goalRoutes.GET("/addNewGoal", goalController.NewGoalSingle)
		goalRoutes.POST("/addNewGoal", goalController.AddGoal)
		// goalRoutes.GET("/:id", goalController.Show)
	}
	categoriesRoutes := router.Group("/categories")
	{
		categoriesRoutes.GET("/listcategories", categoriesController.Index)           // GET request to fetch and display categories at /categories/listcategories
		categoriesRoutes.POST("/listcategories", categoriesController.CreateCategory) // POST request to create a new category at /categories/listcategories

	}
	// }

	return router
}
