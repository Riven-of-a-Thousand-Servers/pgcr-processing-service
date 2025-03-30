package repository

import (
	"database/sql"
	"fmt"
	"pgcr-processing-service/internal/model"
)

type InstanceActivityWeaponStatsRepository struct {
	Conn *sql.DB
}

func (iawsr *InstanceActivityWeaponStatsRepository) AddInstanceWeaponStats(tx *sql.Tx, entity model.InstanceWeaponStats) (*model.InstanceWeaponStats, error) {
	_, err := tx.Exec(`
		INSERT INTO instance_activity_weapon_stats 
		(instance_id, player_character_id, weapon_id, total_kills, total_precision_kills, precision_rates)
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`,
		entity.InstanceId, entity.PlayerCharacterId, entity.WeaponId, entity.TotalKills, entity.TotalPrecisionKills, entity.PrecisionRatio)
	if err != nil {
		return nil, fmt.Errorf(`Error while inserting instance activity weapon stats with arguments: 
			Raid instance: [%d], character Id: [%d], weapon Id: [%d]`,
			entity.InstanceId, entity.PlayerCharacterId, entity.WeaponId)
	}

	return &entity, nil
}
