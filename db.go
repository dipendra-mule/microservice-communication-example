package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// createConnection establishes and returns a Postgres DB connection
func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Fatal("POSTGRES_URL not found in environment variables")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	fmt.Println("✅ Successfully connected to database!")
	return db
}

// User model
type User struct {
	ID    int
	Name  string
	Email string
	Phone string
}

// createUser inserts a user into the database
func createUser(db *sql.DB, name, email string) int {
	var id int
	err := db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		name, email,
	).Scan(&id)
	if err != nil {
		log.Fatal("Error creating user:", err)
	}
	return id
}

// getUsers fetches all users
func getUsers(db *sql.DB) []User {
	rows, err := db.Query("SELECT id, name, email, COALESCE(phone, '') FROM users")
	if err != nil {
		log.Fatal("Error fetching users:", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
		users = append(users, u)
	}
	return users
}
