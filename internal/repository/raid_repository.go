package repository

import (
	"database/sql"
	"fmt"
	"log"
	"rivenbot/internal/model"
)

type RaidRepository struct {
	Conn *sql.DB
}

func (r *RaidRepository) Save(entity model.RaidEntity) (result *model.RaidEntity, err error) {
	transaction, err := r.Conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error creating transaction. %v", err)
	}

	defer func() {
		if err != nil {
			_ = transaction.Rollback()
		}
	}()

	_, err = transaction.Exec(`
    INSERT INTO raid (raid_name, raid_difficulty, is_active, release_date)
    VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		entity.RaidName, entity.RaidDifficulty, entity.IsActive, entity.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting into raid table: %v", err)
	}

	row, err := transaction.Exec(`
      INSERT INTO raid_hash (raid_hash, raid_name, raid_difficulty)
      VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, entity.RaidHash, entity.RaidName, entity.RaidDifficulty)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting raid hash table [%d]: %v", entity.RaidHash, err)
	}

	rows, err := row.RowsAffected()
	if err != nil && rows > 0 {
		log.Printf("Inserted %d new raid hashes for raid [%s:%s]", rows, entity.RaidName, entity.RaidDifficulty)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error while committing raid and raid_hash transaction: %v", err)
	}

	return &entity, nil
}
