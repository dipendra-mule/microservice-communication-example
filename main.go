package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	database := createConnection()
	defer database.Close()

	// Run migrations automatically
	if err := goose.Up(database, "migrations"); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	fmt.Println("Database migrations applied successfully!")
	fmt.Println("Database ready. Starting app...")

	// Example: create a new user
	newUserID := createUser(database, "Alice", "alice1@example.com")
	fmt.Println("Inserted user with ID:", newUserID)

	// Example: list all users
	users := getUsers(database)
	for _, u := range users {
		fmt.Printf("👤 %s (%s)\n", u.Name, u.Email)
	}
}
