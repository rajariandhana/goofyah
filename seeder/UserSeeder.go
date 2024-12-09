package seeder

import (
	"goofyah/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedUser() {
	users := []models.User{
		{Name: "Harry Potter", Email: "harry@gmail.com", Password: "p"},
		{Name: "Hermione Granger", Email: "hermione@gmail.com", Password: "p"},
		{Name: "Ron Weasley", Email: "ron@gmail.com", Password: "p"},
	}
	for _, user := range users {
		existingUser, err := models.GetUserByEmail(user.Email)
		if existingUser != nil || err == nil {
			log.Printf("User with email %s already exists. Skipping...\n", user.Email)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed when hashing password\n")
		}
		user.Password = string(hash)
		if err := models.StoreUser(user); err != nil {
			log.Printf("Error seeding user %s: %v", user.Name, err)
		} else {
			log.Printf("User %s seeded successfully.\n", user.Name)
		}

	}

}
