package service

import (
	"cmp"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"testing"

	"pgcr-processing-service/internal/model"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) AddPlayer(tx *sql.Tx, entity model.PlayerEntity) (*model.PlayerEntity, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*model.PlayerEntity), args.Error(1)
	}
}

func TestProcessPlayers_Success(t *testing.T) {
	// given: An SQL transaction and a processed pgcr
	ppgcr := PPgcrWithGlobalDisplayName()
	tx := sql.Tx{}

	mockedRepository := new(MockPlayerRepository)
	firstPlayer := PgcrToPlayerEntity(ppgcr.PlayerInformation[0])
	secondPlayer := PgcrToPlayerEntity(ppgcr.PlayerInformation[1])

	// Mock repository interactions
	mockedRepository.On("AddPlayer", &tx, firstPlayer).
		Return(&firstPlayer, nil)
	mockedRepository.On("AddPlayer", &tx, secondPlayer).
		Return(&secondPlayer, nil)
	sut := PlayerService{
		PlayerRepository: mockedRepository,
	}

	// when: ProcessPlayers is called
	err := sut.ProcessPlayers(&tx, &ppgcr)

	// then: no error is returned
	if err != nil {
		t.Fatalf("Was not expecting an error, got: %v", err)
	}
}

func TestProcessPlayers_ErrorWhileSavingPlayer(t *testing.T) {
	// given: an SQL transaction and a processed pgcr
	ppgcr := types.ProcessedPostGameCarnageReport{
		PlayerInformation: []types.PlayerData{
			{},
		},
	}
	tx := sql.Tx{}

	mockedRepository := new(MockPlayerRepository)

	// Mock repository interactions
	mockedRepository.On("AddPlayer", mock.Anything, mock.Anything).
		Return(nil, errors.New("Error while inserting player into DB"))

	sut := PlayerService{
		PlayerRepository: mockedRepository,
	}

	// when: ProcessPlayers is called
	err := sut.ProcessPlayers(&tx, &ppgcr)

	// then: An error is returned from the repo
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestPgcrToPlayerEntity_Success_WithGlobalDisplayName(t *testing.T) {
	// given: a ProcessedPgcr
	ppgcr := PPgcrWithGlobalDisplayName()

	// when: PgcrToPlayerEntity is called
	result := PgcrToPlayerEntity(ppgcr.PlayerInformation[0])

	slices.SortFunc(result.Characters, func(a, b model.PlayerCharacterEntity) int {
		return cmp.Compare(a.CharacterId, b.CharacterId)
	})

	// then: the result has correct fields
	assert := assert.New(t)
	firstPlayer := ppgcr.PlayerInformation[0]
	assert.Equal(firstPlayer.DisplayName, result.DisplayName, "Display names should match")
	assert.Equal(firstPlayer.GlobalDisplayName, result.DisplayName, "Global display name should also match display name")
	assert.Equal(int32(firstPlayer.GlobalDisplayNameCode), result.DisplayNameCode, "Global display name code should match the display code")
	assert.Equal(int32(firstPlayer.MembershipType), result.MembershipType, fmt.Sprintf("Expected: [%d], Result: [%d]", firstPlayer.MembershipType, result.MembershipType))
	assert.Equal(firstPlayer.MembershipId, result.MembershipId, "Membership IDs should match")
	assert.Equal(len(firstPlayer.PlayerCharacterInformation), len(result.Characters), "Lengths of characters match")
}

func PPgcrWithGlobalDisplayName() types.ProcessedPostGameCarnageReport {
	return types.ProcessedPostGameCarnageReport{
		PlayerInformation: []types.PlayerData{
			{
				MembershipId:          4611686018440744095,
				MembershipType:        1,
				DisplayName:           "Ceriumz",
				GlobalDisplayName:     "Ceriumz",
				GlobalDisplayNameCode: 1527,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "HUNTER",
					},
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "WARLOCK",
					},
				},
			},
			{
				MembershipId:          4611686018471276567,
				MembershipType:        3,
				DisplayName:           "DUMJ01",
				GlobalDisplayName:     "DUMJ01",
				GlobalDisplayNameCode: 5200,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:     2305843009306384395,
						CharacterEmblem: 3992231368,
						CharacterClass:  "TITAN",
					},
				},
			},
		},
	}
}

func PPgcrWithoutGlobalDisplayName() types.ProcessedPostGameCarnageReport {
	return types.ProcessedPostGameCarnageReport{
		PlayerInformation: []types.PlayerData{
			{
				MembershipId:          4611686018440744095,
				MembershipType:        1,
				DisplayName:           "Ceriumz",
				GlobalDisplayName:     "",
				GlobalDisplayNameCode: 0,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "HUNTER",
					},
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "WARLOCK",
					},
				},
			},
			{
				MembershipId:          4611686018471276567,
				MembershipType:        3,
				DisplayName:           "DUMJ01",
				GlobalDisplayName:     "",
				GlobalDisplayNameCode: 0,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:     2305843009306384395,
						CharacterEmblem: 3992231368,
						CharacterClass:  "TITAN",
					},
				},
			},
		},
	}
}
