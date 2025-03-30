package repository

import (
	"fmt"
	"pgcr-processing-service/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddPlayer_Success(t *testing.T) {
	// given: a player to save
	membershipId := 123838123129
	characters := []model.PlayerCharacterEntity{
		{
			CharacterId:        12775000,
			CharacterClass:     "HUNTER",
			CharacterEmblem:    293012123,
			PlayerMembershipId: int64(membershipId),
		},
		{
			CharacterId:        12775001,
			CharacterClass:     "WARLOCK",
			CharacterEmblem:    293012123,
			PlayerMembershipId: int64(membershipId),
		},
	}
	player := model.PlayerEntity{
		MembershipId:    int64(membershipId),
		DisplayName:     "Deaht",
		DisplayNameCode: 6789,
		MembershipType:  1,
		LastSeen:        time.Now(),
		Characters:      characters,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	repository := PlayerRepositoryImpl{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Error while initializing transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO player").
		WithArgs(player.MembershipId, player.MembershipType, player.DisplayName, player.DisplayNameCode, player.DisplayName, player.LastSeen).
		WillReturnResult(sqlmock.NewResult(1, 1))
	for _, character := range player.Characters {
		mock.ExpectExec("INSERT INTO player_character").
			WithArgs(character.CharacterId, character.CharacterClass, character.CharacterEmblem, character.PlayerMembershipId).
			WillReturnResult(sqlmock.NewResult(2, 2))
	}

	// when: save is called
	result, err := repository.AddPlayer(tx, player)
	if err != nil {
		t.Errorf("Something went wrong when saving player to database: %v", err)
	}

	// then: the returned result is the same entity that was saved
	assert := assert.New(t)
	assert.NotNil(result)
	assert.Equal(player.MembershipId, result.MembershipId, "MembershipId should be the same")
	assert.Equal(player.MembershipType, result.MembershipType, "MembershipType should be the same")
	assert.Equal(player.DisplayName, result.DisplayName, "Display names should be the same")
	assert.Equal(player.DisplayNameCode, result.DisplayNameCode, "Display name codes should be the same")
	assert.NotEmpty(result.Characters, "Characters array should not be empty")

	// and: all characters are the same as the entity passed as an argument
	for i, resultCharacter := range result.Characters {
		assert.Equal(resultCharacter.CharacterId, player.Characters[i].CharacterId, fmt.Sprintf("Character Id should be the same for character at index [%d]", i))
		assert.Equal(resultCharacter.CharacterClass, player.Characters[i].CharacterClass, fmt.Sprintf("Character class should be the same for character at index [%d]", i))
		assert.Equal(resultCharacter.CharacterEmblem, player.Characters[i].CharacterEmblem, fmt.Sprintf("Character emblem should be the same for character at index [%d]", i))
		assert.Equal(resultCharacter.PlayerMembershipId, player.Characters[i].PlayerMembershipId, fmt.Sprintf("Character player membership Id should be the same for characer at index [%d]", i))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database interactions were not met: %v", err)
	}
}

func TestAddPlayer_SuccessOnNoCharactersInserted(t *testing.T) {
	// given: A player to save
	membershipId := 123838123129
	player := model.PlayerEntity{
		MembershipId:    int64(membershipId),
		DisplayName:     "Deaht",
		DisplayNameCode: 6789,
		MembershipType:  1,
		LastSeen:        time.Now(),
		Characters:      []model.PlayerCharacterEntity{},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	repository := PlayerRepositoryImpl{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO player").
		WithArgs(player.MembershipId, player.MembershipType, player.DisplayName, player.DisplayNameCode, player.DisplayName, player.LastSeen).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// when: save is called
	result, err := repository.AddPlayer(tx, player)

	// then: no errors are expected
	if err != nil {
		t.Error("Expecting error, found none")
	}

	// and: all expectations on the db mock should be met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database interactions were not met: %v", err)
	}

	// and: the result should have an empty characters array
	assert := assert.New(t)
	assert.NotNil(result, "Result shouldn't be nil")
	assert.Empty(result.Characters, "Character array is empty for result")
}

func TestAddPlayer_ErrorOnPlayerInsert(t *testing.T) {
	// given: a player to save
	membershipId := 123838123129
	characters := []model.PlayerCharacterEntity{
		{
			CharacterId:        12775000,
			CharacterClass:     "HUNTER",
			CharacterEmblem:    889588491,
			PlayerMembershipId: int64(membershipId),
		},
		{
			CharacterId:        12775001,
			CharacterClass:     "WARLOCK",
			CharacterEmblem:    889588491,
			PlayerMembershipId: int64(membershipId),
		},
	}
	player := model.PlayerEntity{
		MembershipId:    int64(membershipId),
		DisplayName:     "Deaht",
		DisplayNameCode: 6789,
		MembershipType:  1,
		LastSeen:        time.Now(),
		Characters:      characters,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction")
	}

	defer db.Close()

	repository := PlayerRepositoryImpl{
		Conn: db,
	}

	mock.ExpectExec("INSERT INTO player").
		WithArgs(player.MembershipId, player.MembershipType, player.DisplayName, player.DisplayNameCode, player.DisplayName, player.LastSeen).
		WillReturnError(fmt.Errorf("Error inserting player"))

	// when: save is called
	_, err = repository.AddPlayer(tx, player)

	// then: and error is expected
	if err == nil {
		t.Error("Expecting error, found none")
	}

	// and: a rollback and an error in the player insert are expected
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database interactions were not met: %v", err)
	}
}

func TestAddPlayer_ErrorOnPlayerCharacterInsert(t *testing.T) {
	// given: A player to save
	membershipId := 123838123129
	characters := []model.PlayerCharacterEntity{
		{
			CharacterId:        12775000,
			CharacterClass:     "HUNTER",
			CharacterEmblem:    8491823912312,
			PlayerMembershipId: int64(membershipId),
		},
		{
			CharacterId:        12775001,
			CharacterClass:     "WARLOCK",
			CharacterEmblem:    8491823912312,
			PlayerMembershipId: int64(membershipId),
		},
	}
	player := model.PlayerEntity{
		MembershipId:    int64(membershipId),
		DisplayName:     "Deaht",
		DisplayNameCode: 6789,
		MembershipType:  1,
		LastSeen:        time.Now(),
		Characters:      characters,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	repository := PlayerRepositoryImpl{
		Conn: db,
	}

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	mock.ExpectExec("INSERT INTO player").
		WithArgs(player.MembershipId, player.MembershipType, player.DisplayName, player.DisplayNameCode, player.DisplayName, player.LastSeen).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO player_character").
		WithArgs(player.Characters[0].CharacterId, player.Characters[0].CharacterClass, player.Characters[0].CharacterEmblem, player.Characters[0].PlayerMembershipId).
		WillReturnError(fmt.Errorf("Error inserting player character into table"))

	// when: save is called
	_, err = repository.AddPlayer(tx, player)

	// then: no errors are expected
	if err == nil {
		t.Error("Expecting error, found none")
	}

	// and: all expectations on the db mock should be met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database interactions were not met: %v", err)
	}
}
