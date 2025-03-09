package repository

import (
	"fmt"
	"rivenbot/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveRaid(t *testing.T) {
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

	raidRepository := RaidRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO raid_hash`).
		WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// when: save is called for a Raid
	result, err := raidRepository.Save(raid)

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

func TestSaveRaidWithoutRaidHash(t *testing.T) {
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

	raidRepository := RaidRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO raid_hash`).
		WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// when: save is called for a Raid
	result, err := raidRepository.Save(raid)

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

func TestSaveRaidShouldRollbackOnInsertFailure(t *testing.T) {
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

	raidRepository := RaidRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO raid \(raid_name, raid\_difficulty, is\_active, release\_date\)`).
		WithArgs(raid.RaidName, raid.RaidDifficulty, raid.IsActive, raid.ReleaseDate).
		WillReturnError(fmt.Errorf("Some error when inserting into database"))
	mock.ExpectRollback()

	// when: save is called for a Raid
	_, err = raidRepository.Save(raid)
	if err == nil {
		t.Errorf("Was expecting error, got none: %v", err)
	}

	// then: an error should be returned when inserting to raid table fails
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}

func TestSaveRaidShouldRollbackOnExistsQueryFailure(t *testing.T) {
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

	raidRepository := RaidRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO raid_hash").WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnError(fmt.Errorf("Error inserting into raid_hash"))
	mock.ExpectRollback()

	// when: save is called for a Raid
	_, err = raidRepository.Save(raid)
	if err == nil {
		t.Errorf("Was expecting error, got none: %v", err)
	}

	// then: an error is returned when trying to find raid_hash fails
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}

func TestSaveRaidShouldRollbackOnTxCommitFailure(t *testing.T) {
	// given a raid entity from a PGCR
	raid := model.RaidEntity{
		RaidName:       "Last Wish",
		RaidDifficulty: "Normal",
		IsActive:       true,
		ReleaseDate:    time.Now(),
		RaidHash:       4891237913,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error establishing stub connection to database")
	}

	raidRepository := RaidRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO raid").
		WithArgs(raid.RaidName, raid.RaidDifficulty, raid.IsActive, raid.ReleaseDate).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO raid_hash").
		WithArgs(raid.RaidHash, raid.RaidName, raid.RaidDifficulty).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().
		WillReturnError(fmt.Errorf("Some error when commiting a transaction"))

	// when: save is called for a Raid
	_, err = raidRepository.Save(raid)

	// then: an error is expected
	if err == nil {
		t.Errorf("Was expecting error, got none: %v", err)
	}

	// and: the result didn't return an error
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}
