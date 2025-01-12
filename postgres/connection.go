package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Connect to Postgres Database with the required parameters
func Connect() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment file")
	}

	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	if username == "" || password == "" {
		error := fmt.Errorf("Unable to load required parameters to connect to Postgres")
		return nil, error
	}

	connectionString := fmt.Sprintf("username=%s password=&%s dbname=rivenbot sslmode=disabled", username, password)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(time.Second * 30)

	return db, nil
}

// Close all database connections
func Cleanup(db *sql.DB) {
	db.Close()
}
