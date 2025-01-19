package repository

import (
	"database/sql"
	"fmt"
	"log"
	"rivenbot/internal/model"
)

type PgcrRepository struct {
	Conn *sql.DB
}

func (r *PgcrRepository) save(entity model.RaidPgcr) (*model.RaidPgcr, error) {
	transaction, err := r.Conn.Begin()
	if err != nil {
		log.Panic("Error creating transaction")
		return nil, fmt.Errorf("Error creating transaction. %v", err)
	}

	defer transaction.Rollback()

	_, err = transaction.Exec(`INSERT INTO raid_pgcr (instance_id, blob) VALUES ($1, $2)`, entity.InstanceId, entity.Blob)
	if err != nil {
		log.Panicf("Error while inserting into raid_pgcr table: %v", err)
		return nil, err
	}

	err = transaction.Commit()
	if err != nil {
		log.Panicf("Error while commit raid_pgcr transaction for pgcr [%d]: %v", entity.InstanceId, err)
		return nil, err
	}
	return &entity, nil
}
