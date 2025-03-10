package repository

import (
	"database/sql"
	"fmt"
	"rivenbot/internal/model"
)

type InstanceActivityRepository struct {
	Conn *sql.DB
}

func (r *InstanceActivityRepository) AddInstanceActivity(tx *sql.Tx, entity model.InstanceActivityEntity) (*model.InstanceActivityEntity, error) {
	_, err := tx.Exec(`
    INSERT INTO instance_activity_stats
    (
      instance_id, 
      player_membership_id, 
      player_character_id, 
      character_emblem, 
      is_completed,
      kills, 
      deaths, 
      assists, 
      kills_deaths_assists, 
      kills_deaths_ratio, 
      efficiency, 
      duration_seconds, 
      time_played_seconds
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13   
    )`, entity.InstanceId, entity.PlayerMembershipId, entity.PlayerCharacterId,
		entity.CharacterEmblem, entity.IsCompleted, entity.Kills, entity.Deaths, entity.Assists,
		entity.KillsDeathsAssists, entity.KillsDeathsRatio, entity.Efficiency,
		entity.DurationSeconds, entity.TimeplayedSeconds)

	if err != nil {
		return nil, fmt.Errorf("Error inserting activity with instance Id [%d] to database", entity.InstanceId)
	}

	return &entity, nil
}
