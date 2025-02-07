package repository

import (
	"database/sql"
	"fmt"
	"rivenbot/internal/model"
)

type PgcrRepository struct {
	Conn *sql.DB
}

func (r *PgcrRepository) Save(entity model.RaidPgcr) (*model.RaidPgcr, error) {
	transaction, err := r.Conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error creating transaction. %v", err)
	}
	defer transaction.Rollback()

	_, err = transaction.Exec(`INSERT INTO raid_pgcr (instance_id, blob) VALUES ($1, $2)`, entity.InstanceId, entity.Blob)
	if err != nil {
		return nil, fmt.Errorf("Error while inserting into raid_pgcr table: %v", err)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error while commit raid_pgcr transaction for pgcr [%d]: %v", entity.InstanceId, err)
	}
	return &entity, nil
}
