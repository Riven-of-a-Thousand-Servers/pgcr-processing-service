package repository

import (
	"fmt"
	"pgcr-processing-service/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddInstanceWeaponStats_Success(t *testing.T) {
	// given: Instance weapon stats entity
	weaponStats := model.InstanceWeaponStats{
		InstanceId:          1,
		PlayerCharacterId:   49901231,
		WeaponId:            19930812,
		TotalKills:          134,
		TotalPrecisionKills: 35,
		PrecisionRatio:      0.26,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	repository := InstanceActivityWeaponStatsRepository{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO instance_activity_weapon_stats").
		WithArgs(weaponStats.InstanceId, weaponStats.PlayerCharacterId, weaponStats.WeaponId,
			weaponStats.TotalKills, weaponStats.TotalPrecisionKills, weaponStats.PrecisionRatio).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// when: AddInstanceWeaponStats is called
	result, err := repository.AddInstanceWeaponStats(tx, weaponStats)

	// then: the result doesn't return an error
	assert := assert.New(t)

	// and: the result matches the passed entity and all mock expectations are met
	assert.Nil(err)
	assert.Equal(weaponStats.InstanceId, result.InstanceId, "Instance Ids match")
	assert.Equal(weaponStats.PlayerCharacterId, result.PlayerCharacterId, "Player characterIds match")
	assert.Equal(weaponStats.WeaponId, result.WeaponId, "WeaponIds match")
	assert.Equal(weaponStats.TotalKills, result.TotalKills, "Total kills match")
	assert.Equal(weaponStats.TotalPrecisionKills, result.TotalPrecisionKills, "Total precision kills match")
	assert.Equal(weaponStats.PrecisionRatio, result.PrecisionRatio, "Precision ratios match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expecations: %v", err)
	}
}

func TestAddInstanceWeaponStats_ErrorOnInsert(t *testing.T) {
	// given: Instance weapon stats entity
	weaponStats := model.InstanceWeaponStats{
		InstanceId:          1,
		PlayerCharacterId:   49901231,
		WeaponId:            19930812,
		TotalKills:          134,
		TotalPrecisionKills: 35,
		PrecisionRatio:      0.26,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	repository := InstanceActivityWeaponStatsRepository{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO instance_activity_weapon_stats").
		WithArgs(weaponStats.InstanceId, weaponStats.PlayerCharacterId, weaponStats.WeaponId,
			weaponStats.TotalKills, weaponStats.TotalPrecisionKills, weaponStats.PrecisionRatio).
		WillReturnError(fmt.Errorf("Error while inserting into db table"))

	// when: AddInstanceWeaponStats is called
	_, err = repository.AddInstanceWeaponStats(tx, weaponStats)

	// then: an error is returned when trying to insert weapon stats
	if err == nil {
		t.Errorf("Was expecting error, got none: %v", err)
	}

	// and: all mock expecations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expecations: %v", err)
	}
}
