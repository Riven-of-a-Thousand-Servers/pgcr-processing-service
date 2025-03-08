package repository

import (
	"database/sql"
	"fmt"
	"rivenbot/internal/model"
	"time"
)

type PlayerRepository struct {
	Conn *sql.DB
}

func (r *PlayerRepository) Save(entity model.PlayerEntity) (result *model.PlayerEntity, err error) {
	transaction, err := r.Conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("Error creating transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = transaction.Rollback()
		}
	}()

	_, err = transaction.Exec(`
    INSERT INTO player (membership_id, membership_type, global_display_name, global_display_name_code, display_name, last_seen)
    VALUES ($1, $2, $3, $4, $5, $6) 
    ON CONFLICT(membership_id) DO UPDATE
    SET last_seen = EXCLUDED.last_seen`,
		entity.MembershipId, entity.MembershipType, entity.DisplayName, entity.DisplayNameCode, entity.DisplayName, time.Now())

	if err != nil {
		return nil, fmt.Errorf("Error while inserting into player table: %v", err)
	}

	for _, character := range entity.Characters {
		_, err = transaction.Exec(`
          INSERT INTO player_character (character_id, character_class, current_emblem, player_membership_id)
          VALUES ($1, $2, $3, $4) 
          ON CONFLICT (character_id)
          DO UPDATE SET character_emblem = EXCLUDED.current_emblem`,
			character.CharacterId, character.CharacterClass, character.CharacterEmblem, entity.MembershipId)
		if err != nil {
			return nil, fmt.Errorf("Error while upserting player_character [%d]: %v", character.CharacterId, err)
		}
	}

	err = transaction.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error while commiting a player transaction for player [%d]: %v", entity.MembershipId, err)
	}
	return &entity, nil
}
