package main

import (
	"goofyah/config"
	"goofyah/database"
	"goofyah/models"
	"goofyah/routes"
	"goofyah/seeder"
	"log"

	"github.com/joho/godotenv"
)

// var store sessions.Store

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config.LoadConfig()
	db, err := database.Setup()
	models.ConnectDB(db)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	// store := cookie.NewStore([]byte(os.Getenv("SECRET")))
	// store.Options(sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   3600,
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode,
	// })
	seeder.SeedUser(db)
	router := routes.SetupRoutes(db)
	router.Run(":8080")
}
