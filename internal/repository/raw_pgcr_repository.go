package repository

import (
	"database/sql"
	"fmt"
	"rivenbot/internal/model"
)

type RawPgcrRepository struct {
	Conn *sql.DB
}

func (r *RawPgcrRepository) AddRawPgcr(tx *sql.Tx, entity model.RaidPgcr) (*model.RaidPgcr, error) {
	_, err := tx.Exec(`INSERT INTO raid_pgcr (instance_id, blob) VALUES ($1, $2)`, entity.InstanceId, entity.Blob)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting into raid_pgcr table: %v", err)
	}

	return &entity, nil
}
