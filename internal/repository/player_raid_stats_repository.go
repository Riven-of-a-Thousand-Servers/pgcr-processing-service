package repository

import (
	"database/sql"
	"fmt"
	"pgcr-processing-service/internal/model"
)

type PlayerRaidStatsRepository struct {
	Conn *sql.DB
}

func (r *PlayerRaidStatsRepository) AddPlayerRaidStats(tx *sql.Tx, entity model.PlayerRaidStatsEntity) (*model.PlayerRaidStatsEntity, error) {
	_, err := tx.Exec(`
    INSERT INTO player_raid_stats (
    raid_name,
    raid_difficulty,
    player_membership_id,
    kills,
    deaths,
    assists,
    hour_played,
    clears,
    full_clears,
    flawless,
    contest_clear,
    day_one,
    solo,
    duo,
    trio,
    solo_flawless,
    duo_flawless,
    trio_flawless
    )
    VALUES (
      $1, $2, $3, $4, $5, $6, $7, $8, $9,
      $10, $11, $12, $13, $14, $15, $16, $17, $18
    )`, entity.RaidName, entity.RaidDifficulty, entity.PlayerMembershipId, entity.Kills,
		entity.Deaths, entity.Assists, entity.HoursPlayed, entity.Clears, entity.FullClears,
		entity.Flawless, entity.ContestClear, entity.DayOne, entity.Solo, entity.Duo, entity.Trio,
		entity.SoloFlawless, entity.DuoFlawless, entity.TrioFlawless)

	if err != nil {
		return nil, fmt.Errorf("Error while inserting into player table: %v", err)
	}

	return &entity, nil
}
