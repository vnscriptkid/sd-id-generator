package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	connStr := "user=postgres dbname=postgres sslmode=disable password=123456 host=localhost port=5432"

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read and execute init.sql
	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal("Error reading init.sql:", err)
	}

	_, err = db.Exec(string(initSQL))
	if err != nil {
		log.Fatal("Error executing init.sql:", err)
	}

	fmt.Println("Database initialized successfully")

	// Call the id_generator function
	var newID int64
	err = db.QueryRow("SELECT public.id_generator()").Scan(&newID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated ID: %d\n", newID)
}
