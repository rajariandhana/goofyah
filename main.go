package main

import (
	"fmt"
	"goofyah/config"
	"goofyah/database"
	"goofyah/routes"
)

func main() {
	config.LoadConfig()
	db, err := database.Connect()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	router := routes.SetupRoutes(db)
	router.Run(":8080")
}
