package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"rivenbot/internal/model"
	"testing"
)

func TestAddRawRaidPgcr_Success(t *testing.T) {
	// given: a raid pgcr to save
	pgcr := model.RaidPgcr{
		InstanceId: 12377100310231,
		Blob:       []byte{},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	pgcrRepository := RawPgcrRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	// when: save is called
	mock.ExpectExec("INSERT INTO raid_pgcr").WithArgs(pgcr.InstanceId, pgcr.Blob).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := pgcrRepository.AddRawPgcr(tx, pgcr)

	// then: the result didn't return an error
	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(pgcr.InstanceId, result.InstanceId, "Raid PGCR instanceIds should match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddRawRaidPgcr_ErrorOnRaidPgcrInsert(t *testing.T) {
	// given: a raid pgcr to save
	pgcr := model.RaidPgcr{
		InstanceId: 123389859102,
		Blob:       []byte{},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	pgcrRepository := RawPgcrRepository{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	// when: save is called
	mock.ExpectExec("INSERT INTO raid_pgcr").
		WithArgs(pgcr.InstanceId, pgcr.Blob).
		WillReturnError(fmt.Errorf("Some error when inserting into database"))

	// then: we expect an error to be raised when saving
	if _, err = pgcrRepository.AddRawPgcr(tx, pgcr); err == nil {
		t.Errorf("Was expecting error, got none")
	}

	// and: we expect a rollback to be done
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
