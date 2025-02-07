package repository

import (
	"database/sql"
	"fmt"
	"rivenbot/internal/model"
)

type RaidRepository struct {
	Conn *sql.DB
}

func (r *RaidRepository) Save(entity model.RaidEntity) (*model.RaidEntity, error) {
	transaction, err := r.Conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error creating transaction. %v", err)
	}

	defer transaction.Rollback()

	_, err = transaction.Exec(`INSERT INTO raid (raid_name, raid_difficulty, is_active, release_date)
      VALUES ($1, $2, $3, $4)`, entity.RaidName, entity.RaidDifficulty,
		entity.IsActive, entity.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting into raid table: %v", err)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error while committing to raid table: %v", err)
	}

	return &entity, nil
}
