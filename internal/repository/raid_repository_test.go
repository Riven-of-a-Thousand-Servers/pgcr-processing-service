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
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}

func TestSaveRaidShouldRollback(t *testing.T) {
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
	mock.ExpectExec("INSERT INTO raid").WithArgs(raid.RaidName, raid.RaidDifficulty,
		raid.IsActive, raid.ReleaseDate).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("Some error when inserting into database")))
	mock.ExpectCommit()

	// when: save is called for a Raid
	_, err = raidRepository.Save(raid)
	if err != nil {
		t.Errorf("Was expecting error, got none")
	}

	// then: the result didn't return an error
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectation: %v", err)
	}
}
