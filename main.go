package main

import (
	"goofyah/config"
	"goofyah/database"
	"goofyah/routes"
)

func main() {
	config.LoadConfig() // Load environment variables
	database.Connect()  // Initialize DB connection

	router := routes.SetupRoutes() // Set up application routes

	router.Run(":8080") // Start server on port 8080
}
