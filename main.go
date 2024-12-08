package main

import (
	"goofyah/config"
	"goofyah/database"
	"goofyah/routes"
	"goofyah/seeder"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
)

// var store sessions.Store

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config.LoadConfig()
	db, err := database.Setup()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	store := cookie.NewStore([]byte(os.Getenv("SECRET")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	seeder.SeedUser(db)
	router := routes.SetupRoutes(db, store)
	router.Run(":8080")
}
