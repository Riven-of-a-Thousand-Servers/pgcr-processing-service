package repository

import (
	"fmt"
	"pgcr-processing-service/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddRaid_Success(t *testing.T) {
	// given a raid entity from a PGCR
	raid := model.RaidEntity{
		RaidName:       "Last Wish",
		RaidDifficulty: "Normal",
		IsActive:       true,
		ReleaseDate:    time.Now(),
		RaidHash:       1279871289371,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	raidRepository := RaidRepositoryImpl{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO raid_hash`).
		WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// when: save is called for a Raid
	result, err := raidRepository.AddRaidInfo(tx, raid)

	// then: the result didn't return an error
	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(raid.RaidName, result.RaidName, "Raid names should match")
	assert.Equal(raid.RaidDifficulty, result.RaidDifficulty, "Raid difficulty should match")
	assert.Equal(raid.IsActive, result.IsActive, "Raid isActive should match")
	assert.Equal(raid.ReleaseDate, result.ReleaseDate, "Raid release date should match")
	assert.Equal(raid.RaidHash, result.RaidHash, "Raid hashes should match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}

func TestAddRaidInfo_SuccessNoRaidHashAdded(t *testing.T) {
	// given a raid entity from a PGCR
	raid := model.RaidEntity{
		RaidName:       "Last Wish",
		RaidDifficulty: "Normal",
		IsActive:       true,
		ReleaseDate:    time.Now(),
		RaidHash:       1279871289371,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	defer db.Close()

	raidRepository := RaidRepositoryImpl{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO raid_hash`).
		WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// when: save is called for a Raid
	result, err := raidRepository.AddRaidInfo(tx, raid)

	// then: the result didn't return an error
	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(raid.RaidName, result.RaidName, "Raid names should match")
	assert.Equal(raid.RaidDifficulty, result.RaidDifficulty, "Raid difficulty should match")
	assert.Equal(raid.IsActive, result.IsActive, "Raid isActive should match")
	assert.Equal(raid.ReleaseDate, result.ReleaseDate, "Raid release date should match")
	assert.Equal(raid.RaidHash, result.RaidHash, "Raid hashes should match")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}

func TestAddRaidInfo_ErrorOnRaidInsert(t *testing.T) {
	// given a raid entity from a PGCR
	raid := model.RaidEntity{
		RaidName:       "Last Wish",
		RaidDifficulty: "Normal",
		IsActive:       true,
		ReleaseDate:    time.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	raidRepository := RaidRepositoryImpl{
		Conn: db,
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec(`INSERT INTO raid \(raid_name, raid\_difficulty, is\_active, release\_date\)`).
		WithArgs(raid.RaidName, raid.RaidDifficulty, raid.IsActive, raid.ReleaseDate).
		WillReturnError(fmt.Errorf("Some error when inserting into database"))

	// when: save is called for a Raid
	_, err = raidRepository.AddRaidInfo(tx, raid)
	if err == nil {
		t.Errorf("Was expecting error, got none: %v", err)
	}

	// then: an error should be returned when inserting to raid table fails
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}

func TestAddRaidInfo_ErrorOnRaidHashInsert(t *testing.T) {
	// given a raid entity from a PGCR
	raid := model.RaidEntity{
		RaidName:       "Last Wish",
		RaidDifficulty: "Normal",
		IsActive:       true,
		ReleaseDate:    time.Now(),
		RaidHash:       12389172398,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	defer db.Close()

	raidRepository := RaidRepositoryImpl{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO raid_hash").WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnError(fmt.Errorf("Error inserting into raid_hash"))

	// when: save is called for a Raid
	_, err = raidRepository.AddRaidInfo(tx, raid)
	if err == nil {
		t.Errorf("Was expecting error, got none: %v", err)
	}

	// then: an error is returned when trying to find raid_hash fails
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}
