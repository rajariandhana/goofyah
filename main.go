package main

import (
	"goofyah/config"
	"goofyah/database"
	"goofyah/routes"
	"goofyah/seeder"
	"log"

	"github.com/joho/godotenv"
)

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
	seeder.SeedUser()
	router := routes.SetupRoutes(db)
	router.Run("localhost:8080")
}
