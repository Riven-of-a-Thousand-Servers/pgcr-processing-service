package repository

import (
	"fmt"
	"rivenbot/internal/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSavePlayer(t *testing.T) {
	// given: a player to save
	membershipId := 123838123129
	characters := []model.PlayerCharacterEntity{
		{
			CharacterId:        12775000,
			CharacterClass:     "HUNTER",
			CharacterEmblem:    "/some/link/to/bungie",
			PlayerMembershipId: int64(membershipId),
		},
		{
			CharacterId:        12775001,
			CharacterClass:     "WARLOCK",
			CharacterEmblem:    "/some/link/to/bungie",
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

	repository := PlayerRepository{
		Conn: db,
	}

	// when: save is called
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO player").
		WithArgs(player.MembershipId, player.DisplayName, player.DisplayNameCode, player.MembershipType, player.LastSeen).
		WillReturnResult(sqlmock.NewResult(1, 1))
	for _, character := range player.Characters {
		mock.ExpectExec("INSERT INTO player_character").
			WithArgs(character.CharacterClass, character.CharacterEmblem, character.CharacterId, character.PlayerMembershipId).
			WillReturnResult(sqlmock.NewResult(2, 2))
	}

	result, err := repository.Save(player)
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
