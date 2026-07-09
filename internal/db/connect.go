package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Connect to Postgres Database with the required parameters
func Connect(ctx context.Context) (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment file")
	}

	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	if username == "" || password == "" {
		error := fmt.Errorf("Unable to load required parameters to connect to Postgres")
		return nil, error
	}

	url := fmt.Sprintf("postgres://%s:%s@localhost:5432/rivenbot?sslmode=disabled", username, password)
	return OpenDB(url)
}
