package routes

import (
	// "fmt"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("goofyah/views/**/*")

	// router.HTMLRender = loadTemplates("../views")
	// router.Static("/static","./static")

	// tmpl := template.Must(template.New("").ParseGlob("views/templates/*.html"))
	// tmpl = template.Must(tmpl.ParseGlob("views/templates/components/*.html"))
	// router.SetHTMLTemplate(tmpl)

	// var candidates []string
	// candidates = append(candidates, "Anies Baswedan")
	// candidates = append(candidates, "Prabowo Subianto")
	// candidates = append(candidates, "Ganjar Pranowo")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", gin.H{
			"Title": "HOMEPAGE!",
		})
	})
	router.GET("/posts", func(c *gin.Context) {
		c.HTML(200, "posts.html", gin.H{
			"Title": "POSTS page!",
		})
	})

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
