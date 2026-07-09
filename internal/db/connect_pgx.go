package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB(url string) (*sql.DB, error) {
	return sql.Open("pgx", url)
}
