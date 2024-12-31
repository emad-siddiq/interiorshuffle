package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Connect to the default database (e.g., postgres)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to default database: %v\n", err)
	}
	defer db.Close()

	// Run the SQL scripts
	runSQLScript(db, "./db/scripts/create_db.sql")

	// Connect to the created database
	connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the realestate database: %v\n", err)
	}
	defer db.Close()

	// Initialize the database schema and user
	//runSQLScript(db, "./db/scripts/init_tables.sql")
	//runSQLScript(db, "./db/scripts/create_user.sql")
}

func runSQLScript(db *sql.DB, filePath string) {
	script, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read SQL script %s: %v\n", filePath, err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatalf("Failed to execute SQL script %s: %v\n", filePath, err)
	}

	log.Printf("Successfully executed script: %s\n", filePath)
}
