package seeder

import (
	"goofyah/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB) {
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Password: "password"},
		{Name: "Jane Smith", Email: "jane@example.com", Password: "password"},
		{Name: "Bob Johnson", Email: "bob@example.com", Password: "password"},
	}
	for _, user := range users {
		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("Failed when hashing password\n")
			}
			user.Password = string(hash)
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Error seeding user %s: %v", user.Name, err)
			} else {
				log.Printf("User %s seeded successfully.\n", user.Name)
			}
		} else {
			log.Printf("User with email %s already exists. Skipping...\n", user.Email)
		}
	}

}
