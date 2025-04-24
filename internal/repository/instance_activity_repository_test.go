package repository

import (
	"fmt"
	"pgcr-processing-service/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddInstanceActivity_Success(t *testing.T) {
	// given: an instance activity entity
	entity := model.InstanceActivityEntity{
		InstanceId:         2693136604,
		PlayerMembershipId: 4611686018433234646,
		PlayerCharacterId:  2305843009260849623,
		CharacterEmblem:    3115055261,
		IsCompleted:        true,
		Kills:              58,
		Deaths:             4,
		Assists:            31,
		KillsDeathsAssists: 18.38,
		KillsDeathsRatio:   14.50,
		DurationSeconds:    15987,
		TimeplayedSeconds:  1333,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	repository := InstanceActivityRepositoryImpl{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO instance_activity_stats").
		WithArgs(entity.InstanceId, entity.PlayerMembershipId, entity.PlayerCharacterId,
			entity.CharacterEmblem, entity.IsCompleted, entity.Kills, entity.Deaths, entity.Assists,
			entity.KillsDeathsAssists, entity.KillsDeathsRatio,
			entity.DurationSeconds, entity.TimeplayedSeconds).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// when: save is called
	result, err := repository.AddInstanceActivity(tx, entity)

	if err != nil {
		t.Fatalf("An error should not be expected")
	}

	// then: the result is expected to have same fields as the input
	assert := assert.New(t)
	assert.Equal(entity.InstanceId, result.InstanceId, "Instance IDs should be the same")
	assert.Equal(entity.PlayerMembershipId, result.PlayerMembershipId, "PlayerMembershipIds should be the same")
	assert.Equal(entity.PlayerCharacterId, result.PlayerCharacterId, "PlayerCharacterIds should be the same")
	assert.Equal(entity.CharacterEmblem, result.CharacterEmblem, "CharacterEmblems should be the same")
	assert.Equal(entity.IsCompleted, result.IsCompleted, "IsCompleted should be the same")
	assert.Equal(entity.Kills, result.Kills, "Kills should be the same")
	assert.Equal(entity.Deaths, result.Deaths, "Deaths should be the same")
	assert.Equal(entity.Assists, result.Assists, "Assists should be the same")
	assert.Equal(entity.KillsDeathsAssists, result.KillsDeathsAssists, "KDAs should be the same")
	assert.Equal(entity.KillsDeathsRatio, result.KillsDeathsRatio, "KDRs should be the same")
	assert.Equal(entity.DurationSeconds, result.DurationSeconds, "Durations should be the same")
	assert.Equal(entity.TimeplayedSeconds, result.TimeplayedSeconds, "Time played should be the same")

	// and: expected db interactions should be correct
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expecations: %v", err)
	}
}

func TestAddInstanceActivity_ErrorOnActivityInstanceInsert(t *testing.T) {
	// given: an instance activity to save
	entity := model.InstanceActivityEntity{
		InstanceId:         2693136604,
		PlayerMembershipId: 4611686018433234646,
		PlayerCharacterId:  2305843009260849623,
		CharacterEmblem:    3115055261,
		IsCompleted:        true,
		Kills:              58,
		Deaths:             4,
		Assists:            31,
		KillsDeathsAssists: 18.38,
		KillsDeathsRatio:   14.50,
		DurationSeconds:    15987,
		TimeplayedSeconds:  1333,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO instance_activity_stats").
		WithArgs(entity.InstanceId, entity.PlayerMembershipId, entity.PlayerCharacterId,
			entity.CharacterEmblem, entity.IsCompleted, entity.Kills, entity.Deaths, entity.Assists,
			entity.KillsDeathsAssists, entity.KillsDeathsRatio,
			entity.DurationSeconds, entity.TimeplayedSeconds).
		WillReturnError(fmt.Errorf("Something happened while inserting into table"))

	repository := InstanceActivityRepositoryImpl{
		Conn: db,
	}

	// when: save is called
	_, err = repository.AddInstanceActivity(tx, entity)

	// then: an error is expected
	if err == nil {
		t.Fatalf("Expected error, found none")
	}

	// and: db interactions are expected to be correct
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %v", err)
	}
}
