package repository

import (
	"database/sql"
	"rivenbot/internal/model"
)

type PlayerRepository struct {
	Conn *sql.DB
}

func (r *PlayerRepository) Save(entity model.PlayerEntity) (result *model.RaidEntity, err error) {
	// TODO: Finish implementing this repo
	return nil, nil
}
