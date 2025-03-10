package repository

import (
	"database/sql"
)

type Repository[T any] interface {
	save(tx *sql.Tx, entity T) (*T, error)
}
