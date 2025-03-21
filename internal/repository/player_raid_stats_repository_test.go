package repository

import (
	"fmt"
	"rivenbot/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

func TestAddPlayerRaidStats_Success(t *testing.T) {
	// given: a player raid stats entity
	playerStats := model.PlayerRaidStatsEntity{
		RaidName:           types.LAST_WISH,
		RaidDifficulty:     types.NORMAL,
		PlayerMembershipId: 4611686018440744095,
		Kills:              1223,
		Deaths:             27,
		Assists:            174,
		HoursPlayed:        18989,
		Clears:             43,
		FullClears:         38,
		Flawless:           true,
		ContestClear:       true,
		DayOne:             false,
		Solo:               false,
		Duo:                false,
		Trio:               false,
		SoloFlawless:       false,
		DuoFlawless:        false,
		TrioFlawless:       false,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	repository := PlayerRaidStatsRepository{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	mock.ExpectExec("INSERT INTO player_raid_stats").
		WithArgs(playerStats.RaidName, playerStats.RaidDifficulty, playerStats.PlayerMembershipId,
			playerStats.Kills, playerStats.Deaths, playerStats.Assists, playerStats.HoursPlayed,
			playerStats.Clears, playerStats.FullClears, playerStats.Flawless, playerStats.ContestClear,
			playerStats.DayOne, playerStats.Solo, playerStats.Duo, playerStats.Trio, playerStats.SoloFlawless,
			playerStats.DuoFlawless, playerStats.TrioFlawless).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// when: AddPlayerRaidStats is called
	result, err := repository.AddPlayerRaidStats(tx, playerStats)

	if err != nil {
		t.Fatalf("Not expecting error but one was thrown: %v", err)
	}

	assert := assert.New(t)
	assert.Equal(playerStats.RaidName, result.RaidName, "Raid name should match")
	assert.Equal(playerStats.RaidDifficulty, result.RaidDifficulty, "Raid difficulty should match")
	assert.Equal(playerStats.PlayerMembershipId, result.PlayerMembershipId, "Players membership ids should match")
	assert.Equal(playerStats.Kills, result.Kills, "Kills should be the same")
	assert.Equal(playerStats.Deaths, result.Deaths, "Deaths should be the same")
	assert.Equal(playerStats.Assists, result.Assists, "Assits should be the same")
	assert.Equal(playerStats.HoursPlayed, result.HoursPlayed, "Hours played should be the same")
	assert.Equal(playerStats.Clears, result.Clears, "Clears should be the same")
	assert.Equal(playerStats.FullClears, result.FullClears, "Full clears should the same")
	assert.Equal(playerStats.Flawless, result.Flawless, "Flawless should be equal")
	assert.Equal(playerStats.ContestClear, result.ContestClear, "Contest clears should be equal")
	assert.Equal(playerStats.DayOne, result.DayOne, "Day one should be equal")
	assert.Equal(playerStats.Solo, result.Solo, "Solos should be equal")
	assert.Equal(playerStats.Duo, result.Duo, "Duos should be equal")
	assert.Equal(playerStats.Trio, result.Trio, "Trios should be equal")
	assert.Equal(playerStats.SoloFlawless, result.SoloFlawless, "Solo flawless should be equal")
	assert.Equal(playerStats.DuoFlawless, result.DuoFlawless, "Duo flawless should be equal")
	assert.Equal(playerStats.TrioFlawless, result.TrioFlawless, "Trio flawless should be equal")

	// and: db interactions are expected to be correct
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %v", err)
	}
}

func TestAddPlayerRaidStats_ErrorOnPlayerRaidStatsInsert(t *testing.T) {
	// given: a player stats entity to insert
	playerStats := model.PlayerRaidStatsEntity{
		RaidName:           types.LAST_WISH,
		RaidDifficulty:     types.NORMAL,
		PlayerMembershipId: 4611686018440744095,
		Kills:              1223,
		Deaths:             27,
		Assists:            174,
		HoursPlayed:        18989,
		Clears:             43,
		FullClears:         38,
		Flawless:           true,
		ContestClear:       true,
		DayOne:             false,
		Solo:               false,
		Duo:                false,
		Trio:               false,
		SoloFlawless:       false,
		DuoFlawless:        false,
		TrioFlawless:       false,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	repository := PlayerRaidStatsRepository{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	mock.ExpectExec("INSERT INTO player_raid_stats").
		WithArgs(playerStats.RaidName, playerStats.RaidDifficulty, playerStats.PlayerMembershipId,
			playerStats.Kills, playerStats.Deaths, playerStats.Assists, playerStats.HoursPlayed,
			playerStats.Clears, playerStats.FullClears, playerStats.Flawless, playerStats.ContestClear,
			playerStats.DayOne, playerStats.Solo, playerStats.Duo, playerStats.Trio, playerStats.SoloFlawless,
			playerStats.DuoFlawless, playerStats.TrioFlawless).
		WillReturnError(fmt.Errorf("Something happened when inserting player raid stats"))

	// when: AddPlayerRaidStats is called
	_, err = repository.AddPlayerRaidStats(tx, playerStats)

	if err == nil {
		t.Fatalf("Expecting error, found none")
	}

	// and: db interactions are expected to be correct
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %v", err)
	}
}
