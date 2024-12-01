package routes

import (
	// "fmt"

	"fmt"
	"goofyah/controllers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	// router.LoadHTMLGlob("views/*.html")
	// router.LoadHTMLGlob("views/user/*.html")
	// router.LoadHTMLFiles("views/header.html", "views/footer.html", "views/index.html")
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

	router.LoadHTMLFiles(htmlFiles...)

	var list = []string{"anies", "prabowo", "ganjar"}
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
			"list":  list,
		})
	})

	userController := controllers.NewUserController(db)

	router.POST("/users/:id", func(c *gin.Context) {
		if c.DefaultPostForm("_method", "") == "PUT" {
			userController.Update(c)
		} else if c.DefaultPostForm("_method", "") == "DELETE" {
			userController.Destroy(c)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method"})
		}
	})
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", userController.Index)
		userRoutes.GET("/:id", userController.Show)
		// userRoutes.GET("/create", userController.Create)
		userRoutes.POST("/store", userController.Store)
		// userRoutes.GET("/:id", userController.Edit)
		userRoutes.PUT("/update/:id", userController.Update)
		userRoutes.DELETE("/destroy/:id", userController.Destroy)
	}

	// router.HTMLRender = loadTemplates("../views")
	// router.Static("/static","./static")

	// tmpl := template.Must(template.New("").ParseGlob("views/templates/*.html"))
	// tmpl = template.Must(tmpl.ParseGlob("views/templates/components/*.html"))
	// router.SetHTMLTemplate(tmpl)

	// var candidates []string
	// candidates = append(candidates, "Anies Baswedan")
	// candidates = append(candidates, "Prabowo Subianto")
	// candidates = append(candidates, "Ganjar Pranowo")
	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(200, "home.html", gin.H{
	// 		"Title": "HOMEPAGE!",
	// 	})
	// })
	// router.GET("/posts", func(c *gin.Context) {
	// 	c.HTML(200, "posts.html", gin.H{
	// 		"Title": "POSTS page!",
	// 	})
	// })

	// userRoutes := router.Group("/users")
	// {
	// 	userRoutes.GET("/", controllers.GetAllUsers)
	// 	userRoutes.POST("/", controllers.CreateUser)
	// 	userRoutes.GET("/:id", controllers.GetUser)
	// 	userRoutes.PUT("/:id", controllers.UpdateUser)
	// 	userRoutes.DELETE("/:id", controllers.DeleteUser)
	// }
	return router
}
