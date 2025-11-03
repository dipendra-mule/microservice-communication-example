package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // 👈 important: PostgreSQL driver import
)

func main() {
	db := createConnection()
	defer db.Close()

	fmt.Println("Database connection established")

	// create the table if not exists
	createTable(db)

	// create a user
	newUserID := createUser(db, "John Doe", "john.doe@example.com")
	fmt.Printf("New user created with ID: %d\n", newUserID)

	// get all users
	users := getUsers(db)
	for _, u := range users {
		fmt.Printf("User: %s\n  Email: %s\n", u.Name, u.Email)
	}
}

func createConnection() *sql.DB {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get connection URL from env
	dbURL := os.Getenv("POSTGRES_URL")
	if dbURL == "" {
		log.Fatal("POSTGRES_URL not found in environment variables")
	}

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(20),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// _, err := db.Exec(query)
	// if err != nil {
	// 	log.Fatal("Error creating table:", err)
	// }
	// fmt.Println("Table created successfully")
}

type User struct {
	ID    int
	Name  string
	Email string
}

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

// getUsers fetches all users from the database
func getUsers(db *sql.DB) []User {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Fatal("Error fetching users:", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
		users = append(users, u)
	}
	return users
}
