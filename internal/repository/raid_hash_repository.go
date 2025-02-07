package repository

import (
	"database/sql"
	"fmt"
	"rivenbot/internal/model"
)

type RaidHashRepository struct {
	Conn *sql.Conn
}
