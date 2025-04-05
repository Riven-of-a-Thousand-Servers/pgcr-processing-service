package repository

import (
	"database/sql"
	"fmt"
	"log"
	"pgcr-processing-service/internal/model"
)

type RaidRepositoryImpl struct {
	Conn *sql.DB
}

type RaidRepository interface {
	AddRaidInfo(tx *sql.Tx, entity model.RaidEntity) (*model.RaidEntity, error)
}

func (r *RaidRepositoryImpl) AddRaidInfo(tx *sql.Tx, entity model.RaidEntity) (*model.RaidEntity, error) {
	_, err := tx.Exec(`
    INSERT INTO raid (raid_name, raid_difficulty, is_active, release_date)
    VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		entity.RaidName, entity.RaidDifficulty, entity.IsActive, entity.ReleaseDate)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting into raid table: %v", err)
	}

	row, err := tx.Exec(`
      INSERT INTO raid_hash (raid_hash, raid_name, raid_difficulty)
      VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, entity.RaidHash, entity.RaidName, entity.RaidDifficulty)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting raid hash table [%d]: %v", entity.RaidHash, err)
	}

	rows, err := row.RowsAffected()
	if err != nil && rows > 0 {
		log.Printf("Inserted %d new raid hashes for raid [%s:%s]", rows, entity.RaidName, entity.RaidDifficulty)
	}

	return &entity, nil
}
