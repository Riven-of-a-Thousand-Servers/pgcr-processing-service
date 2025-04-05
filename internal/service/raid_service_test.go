package service

import (
	"database/sql"
	"fmt"
	"pgcr-processing-service/internal/model"
	"testing"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRaidRepository struct {
	mock.Mock
}

func (m *MockRaidRepository) AddRaidInfo(tx *sql.Tx, entity model.RaidEntity) (*model.RaidEntity, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*model.RaidEntity), args.Error(1)
	}
}

func TestProcessRaid_Success(t *testing.T) {
	// given: a processed pgcr
	ppgcr := types.ProcessedPostGameCarnageReport{
		RaidName:       types.CROTAS_END,
		RaidDifficulty: types.MASTER,
		ActivityHash:   12381923819231,
	}

	tx := sql.Tx{}
	mockRepository := new(MockRaidRepository)

	sut := RaidService{
		RaidRepository: mockRepository,
	}

	resultEntity := MapPgcrToRaidEntity(&ppgcr)
	mockRepository.On("AddRaidInfo", &tx, resultEntity).
		Return(&resultEntity, nil)

	// when: ProcessRaid is called
	err := sut.ProcessRaid(&tx, &ppgcr)

	// then: no error is thrown
	if err != nil {
		t.Fatalf("Not expecting error, got: %v", err)
	}
}

func TestProcessRaid_ErrorOnRepositoryCall(t *testing.T) {
	// given: a processed pgcr
	ppgcr := types.ProcessedPostGameCarnageReport{
		RaidName:       types.SALVATIONS_EDGE,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1238123102930,
	}

	tx := sql.Tx{}

	mockRepository := new(MockRaidRepository)

	sut := RaidService{
		RaidRepository: mockRepository,
	}

	resultEntity := MapPgcrToRaidEntity(&ppgcr)
	mockRepository.On("AddRaidInfo", &tx, resultEntity).
		Return(nil, fmt.Errorf("Something happened while inserting raid into DB"))

	// when: Process raid is called
	err := sut.ProcessRaid(&tx, &ppgcr)

	// then: an error is returned
	if err == nil {
		t.Fatal("Expecting error but found none")
	}
}

var mappingInputs = map[string]types.ProcessedPostGameCarnageReport{
	"Salvation's Edge mapping test": {
		RaidName:       types.SALVATIONS_EDGE,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Crota's End mapping test": {
		RaidName:       types.CROTAS_END,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Root of Nightmares mapping test": {
		RaidName:       types.ROOT_OF_NIGHTMARES,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"King's Fall mapping test": {
		RaidName:       types.KINGS_FALL,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Vow of the Disciple mapping test": {
		RaidName:       types.VOW_OF_THE_DISCIPLE,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Vault of Glass mapping test": {
		RaidName:       types.VAULT_OF_GLASS,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Deep Stone Crypt mapping test": {
		RaidName:       types.DEEP_STONE_CRYPT,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Garden of Salvation mapping test": {
		RaidName:       types.GARDEN_OF_SALVATION,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Crown of Sorrow mapping test": {
		RaidName:       types.CROWN_OF_SORROW,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Last Wish mapping test": {
		RaidName:       types.LAST_WISH,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Scourge of the Past mapping test": {
		RaidName:       types.SCOURGE_OF_THE_PAST,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Spire of Stars mapping test": {
		RaidName:       types.SPIRE_OF_STARS,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Eater of Worlds mapping test": {
		RaidName:       types.EATER_OF_WORLDS,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
	"Leviathan": {
		RaidName:       types.LEVIATHAN,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1,
	},
}

func TestMapPgcrToRaidEntity(t *testing.T) {
	for testName, input := range mappingInputs {
		t.Run(testName, func(t *testing.T) {
			// when: MapPgcrToRaidEntity is called
			result := MapPgcrToRaidEntity(&input)

			// then: the result is correct
			assert := assert.New(t)
			assert.Equal(result.RaidName, string(input.RaidName), "Raid names should match")
			assert.Equal(result.RaidDifficulty, string(input.RaidDifficulty), "Raid difficulties should match")
			assert.Equal(result.RaidHash, input.ActivityHash, "Raid hashes should match")
			assert.Equal(result.IsActive, activeRaids[input.RaidName], "If a raid is active should match what's stored in memory")
			assert.Equal(result.ReleaseDate, raidReleaseDates[input.RaidName])
		})
	}
}
