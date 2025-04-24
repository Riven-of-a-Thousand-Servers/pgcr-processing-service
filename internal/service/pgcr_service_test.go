package service

import (
	"cmp"
	"database/sql"
	"errors"
	"fmt"
	"pgcr-processing-service/internal/model"
	"slices"
	"testing"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddRaidsToTx_Success(t *testing.T) {
	// given: db transaction and a processed pgcr
	ppgcr := types.ProcessedPostGameCarnageReport{
		RaidName:       types.CROTAS_END,
		RaidDifficulty: types.MASTER,
		ActivityHash:   12381923819231,
	}

	tx := sql.Tx{}

	mockRaidRepository := new(MockRaidRepository)
	sut := PgcrServiceImpl{
		RaidRepository: mockRaidRepository,
	}

	resultEntity := MapPgcrToRaidEntity(&ppgcr)
	mockRaidRepository.On("AddRaidInfo", &tx, resultEntity).
		Return(&resultEntity, nil)

	// when: ProcessRaid is called
	err := sut.addRaidInfoToTx(&tx, &ppgcr)

	// then: no error is thrown
	if err != nil {
		t.Fatalf("Not expecting error, got: %v", err)
	}
}

func TestAddRaidsToTx_ErrorOnRepositoryCall(t *testing.T) {
	// given: a processed pgcr
	ppgcr := types.ProcessedPostGameCarnageReport{
		RaidName:       types.SALVATIONS_EDGE,
		RaidDifficulty: types.NORMAL,
		ActivityHash:   1238123102930,
	}

	tx := sql.Tx{}

	mockRaidRepository := new(MockRaidRepository)
	sut := PgcrServiceImpl{
		RaidRepository: mockRaidRepository,
	}

	resultEntity := MapPgcrToRaidEntity(&ppgcr)
	mockRaidRepository.On("AddRaidInfo", &tx, resultEntity).
		Return(nil, fmt.Errorf("Something happened while inserting raid into DB"))

	// when: Process raid is called
	err := sut.addRaidInfoToTx(&tx, &ppgcr)

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

func TestAddPlayersToTx_Success(t *testing.T) {
	// given: An SQL transaction and a processed pgcr
	ppgcr := PPgcrWithGlobalDisplayName()
	tx := sql.Tx{}

	firstPlayer := MapPlayerEntity(ppgcr.PlayerInformation[0])
	secondPlayer := MapPlayerEntity(ppgcr.PlayerInformation[1])

	// Mock repository interactions
	mockPlayerRepository := new(MockPlayerRepository)
	mockPlayerRepository.On("AddPlayer", &tx, firstPlayer).
		Return(&firstPlayer, nil)
	mockPlayerRepository.On("AddPlayer", &tx, secondPlayer).
		Return(&secondPlayer, nil)

	sut := PgcrServiceImpl{
		PlayerRepository: mockPlayerRepository,
	}

	// when: ProcessPlayers is called
	err := sut.addPlayerInfoToTx(&tx, &ppgcr)

	// then: no error is returned
	if err != nil {
		t.Fatalf("Was not expecting an error, got: %v", err)
	}
}

func TestAddPlayersToTx_ErrorWhileSavingPlayer(t *testing.T) {
	// given: an SQL transaction and a processed pgcr
	ppgcr := types.ProcessedPostGameCarnageReport{
		PlayerInformation: []types.PlayerData{
			{},
		},
	}
	tx := sql.Tx{}

	// Mock repository interactions
	mockPlayerRepository := new(MockPlayerRepository)
	mockPlayerRepository.On("AddPlayer", mock.Anything, mock.Anything).
		Return(nil, errors.New("Error while inserting player into DB"))

	sut := PgcrServiceImpl{
		PlayerRepository: mockPlayerRepository,
	}

	// when: ProcessPlayers is called
	err := sut.addPlayerInfoToTx(&tx, &ppgcr)

	// then: An error is returned from the repo
	if err == nil {
		t.Fatal("Expected an error, got none")
	}
}

func TestMapPlayer_Success_WithGlobalDisplayName(t *testing.T) {
	// given: a ProcessedPgcr
	ppgcr := PPgcrWithGlobalDisplayName()

	// when: MapPlayer is called
	result := MapPlayerEntity(ppgcr.PlayerInformation[0])

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

func TestMapPlayer_Success_WithoutGlobalDisplayName(t *testing.T) {
	// given: some player data
	ppgcr := PPgcrWithoutGlobalDisplayName()

	// when: MapPlayer is called
	result := MapPlayerEntity(ppgcr.PlayerInformation[0])

	slices.SortFunc(result.Characters, func(a, b model.PlayerCharacterEntity) int {
		return cmp.Compare(a.CharacterId, b.CharacterId)
	})

	// then: the result has correct fields
	assert := assert.New(t)
	firstPlayer := ppgcr.PlayerInformation[0]
	assert.Equal(firstPlayer.DisplayName, result.DisplayName, "Display names should match")
	assert.Equal(int32(0), result.DisplayNameCode, "Global display name code should be zero")
	assert.Equal(int32(firstPlayer.MembershipType), result.MembershipType, fmt.Sprintf("Expected: [%d], Result: [%d]", firstPlayer.MembershipType, result.MembershipType))
	assert.Equal(firstPlayer.MembershipId, result.MembershipId, "Membership IDs should match")
	assert.Equal(len(firstPlayer.PlayerCharacterInformation), len(result.Characters), "Lengths of characters match")
}

func TestMapInstanceActivityEntity_Success(t *testing.T) {
	// given: processed pgcr, player data and character data
	ppgcr := types.ProcessedPostGameCarnageReport{
		InstanceId: 12390,
	}

	playerData := types.PlayerData{
		MembershipId: 123,
	}

	characterData := types.PlayerCharacterInformation{
		CharacterId:       1234,
		CharacterEmblem:   12345,
		ActivityCompleted: false,
		Kills:             53,
		Deaths:            1,
		Assists:           3,
		AbilityInformation: types.CharacterAbilityInformation{
			SuperKills:   3,
			GrenadeKills: 5,
			MeleeKills:   19,
		},
		Kdr:               3.1,
		Kda:               4.5,
		TimePlayedSeconds: 341,
	}

	// when: mapInstanceActvityEntity is called
	result := MapInstanceActivityEntity(&ppgcr, &playerData, &characterData)

	// then: the fields in the result are correct
	assert := assert.New(t)
	assert.Equal(ppgcr.InstanceId, result.InstanceId, "Instance IDs should be the same")
	assert.Equal(playerData.MembershipId, result.PlayerMembershipId, "Player membership Ids should match")
	assert.Equal(characterData.CharacterId, result.PlayerCharacterId, "Character Id should match")
	assert.Equal(characterData.CharacterEmblem, result.CharacterEmblem, "Emblems should match")
	assert.Equal(characterData.ActivityCompleted, result.IsCompleted, "Is completed should match activity completed")
	assert.Equal(int32(characterData.Kills), result.Kills, "Kills should match")
	assert.Equal(int32(characterData.Deaths), result.Deaths, "Deaths should match")
	assert.Equal(int32(characterData.Assists), result.Assists, "Assists should match")
	assert.Equal(characterData.AbilityInformation.SuperKills, result.SuperKills, "Super kills should match")
	assert.Equal(characterData.AbilityInformation.GrenadeKills, result.GrenadeKills, "Grenade kills should match")
	assert.Equal(characterData.AbilityInformation.MeleeKills, result.MeleeKills, "Melee kills should match")
	assert.Equal(characterData.Kda, result.KillsDeathsAssists, "KDA should match")
	assert.Equal(characterData.Kdr, result.KillsDeathsRatio, "KDR should match")
	assert.Equal(int64(ppgcr.EndTime.Sub(ppgcr.StartTime)), result.DurationSeconds, "Duration should match the subtraction of end - start")
	assert.Equal(int64(characterData.TimePlayedSeconds), result.TimeplayedSeconds, "Time played should match")
}

func TestAddInstanceInfoToTx_Success(t *testing.T) {
	// given: sql transaction and processed pgcr
	tx := sql.Tx{}
	ppgcr := FlawlessDuo()

	mockInstanceRepository := new(MockInstanceActivityRepository)
	sut := PgcrServiceImpl{
		InstanceActivityRepository: mockInstanceRepository,
	}

	player1 := MapInstanceActivityEntity(&ppgcr, &ppgcr.PlayerInformation[0], &ppgcr.PlayerInformation[0].PlayerCharacterInformation[0])
	player2 := MapInstanceActivityEntity(&ppgcr, &ppgcr.PlayerInformation[1], &ppgcr.PlayerInformation[1].PlayerCharacterInformation[0])

	mockInstanceRepository.On("AddInstanceActivity", &tx, player1).
		Return(&player1, nil)
	mockInstanceRepository.On("AddInstanceActivity", &tx, player2).
		Return(&player2, nil)

	// when: addInstanceInfoToTx is called
	err := sut.addInstanceInfoToTx(&tx, &ppgcr)

	// then: no error is thrown
	if err != nil {
		t.Fatalf("Not expecting error, got: %v", err)
	}

	if !mockInstanceRepository.AssertExpectations(t) {
		t.Fatalf("Some expectations were not met")
	}
}

func TestAddInstanceInfoToTx_ErrorRepositoryCall(t *testing.T) {
	// given: sql transaction and processed pgcr
	tx := sql.Tx{}
	ppgcr := FlawlessDuo()

	entity := MapInstanceActivityEntity(&ppgcr, &ppgcr.PlayerInformation[0], &ppgcr.PlayerInformation[0].PlayerCharacterInformation[0])

	mockInstanceRepository := new(MockInstanceActivityRepository)
	mockInstanceRepository.On("AddInstanceActivity", &tx, entity).
		Return(nil, errors.New("Something happened while adding entity"))

	sut := PgcrServiceImpl{
		InstanceActivityRepository: mockInstanceRepository,
	}

	// when: addInstanceInfoToTx is called
	err := sut.addInstanceInfoToTx(&tx, &ppgcr)

	// then: an error is returned
	if err == nil {
		t.Fatal("Expecting error, found none")
	}

	if !mockInstanceRepository.AssertExpectations(t) {
		t.Fatal("Some expectations were not met")
	}
}

func TestAddWeaponInfoToTx_Success(t *testing.T) {
	// given: an SQL transaction and a processed post game carnage report
	tx := sql.Tx{}
	ppgcr := PpgcrWithWeaponStats()

	mockRedisService := new(MockRedisService)
	mockWeaponRepository := new(MockWeaponRepository)
	mockActivityWeaponStatsRepository := new(MockActivityWeaponStatsRepository)

	mockRedisService.On("GetManifestEntity", mock.Anything, "37189237").
		Return(&types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Icon: "some/icon/url/2",
				Name: "Izanagi's Burden",
			},
		}, nil)

	mockRedisService.On("GetManifestEntity", mock.Anything, "8978912").
		Return(&types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Icon: "some/icon/url/1",
				Name: "Fatebringer (Timelost)",
			},
		}, nil)

	mockWeaponRepository.On("AddWeapon", &tx, mock.Anything).
		Return(&model.WeaponEntity{}, nil)

	mockActivityWeaponStatsRepository.On("AddInstanceWeaponStats", &tx, mock.Anything).
		Return(&model.InstanceWeaponStats{}, nil)

	sut := PgcrServiceImpl{
		Redis:                         mockRedisService,
		WeaponRepository:              mockWeaponRepository,
		InstanceWeaponStatsRepository: mockActivityWeaponStatsRepository,
	}

	// when: addWeaponInfoToTx is called
	err := sut.addWeaponInfoToTx(&tx, &ppgcr)

	// then: no error is returned
	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}
}

func TestAddWeaponInfoToTx_ErrorWhileCallingManifest(t *testing.T) {
	// given: an SQL transaction and a processed pgcr
	ppgcr := PpgcrWithWeaponStats()
	tx := sql.Tx{}

	mockRedisService := new(MockRedisService)

	mockRedisService.On("GetManifestEntity", mock.Anything, "37189237").
		Return(nil, fmt.Errorf("Some error getting manifest entity"))

	sut := PgcrServiceImpl{
		Redis: mockRedisService,
	}

	// when: addWeaponInfoToTx is called
	err := sut.addWeaponInfoToTx(&tx, &ppgcr)

	// then: an error is returned when calling the manifest
	if err == nil {
		t.Fatalf("Expecting error, found none")
	}
}

func TestAddWeaponInfoToTx_ErrorWhileSavingWeapon(t *testing.T) {
	// given: an SQL transaction and a processed pgcr
	ppgcr := PpgcrWithWeaponStats()
	tx := sql.Tx{}

	mockRedisService := new(MockRedisService)
	mockWeaponRepository := new(MockWeaponRepository)

	mockRedisService.On("GetManifestEntity", mock.Anything, "37189237").
		Return(&types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Icon: "some/icon/url/2",
				Name: "Izanagi's Burden",
			},
			EquippingBlock: types.EquippingBlock{
				AmmoType:              1,
				EquipmentSlotTypeHash: 1498876634,
			},
		}, nil)

	weapon := model.WeaponEntity{
		WeaponHash:          37189237,
		WeaponIcon:          "some/icon/url/2",
		WeaponName:          "Izanagi's Burden",
		WeaponDamageType:    "Kinetic",
		WeaponEquipmentSlot: "Primary",
	}
	mockWeaponRepository.On("AddWeapon", &tx, weapon).
		Return(nil, fmt.Errorf("Error adding weapon"))

	sut := PgcrServiceImpl{
		Redis:            mockRedisService,
		WeaponRepository: mockWeaponRepository,
	}

	// when: addWeaponInfoToTx is called
	err := sut.addWeaponInfoToTx(&tx, &ppgcr)

	// then: an error is returned when saving weapon
	if err == nil {
		t.Fatalf("Expecting error, found none")
	}
}

func TestAddWeaponInfoToTx_ErrorSavingInstanceWeaponStats(t *testing.T) {
	// given: an SQL transaction and a ppgcr
	ppgcr := PpgcrWithWeaponStats()
	tx := sql.Tx{}

	mockRedisService := new(MockRedisService)
	mockWeaponRepository := new(MockWeaponRepository)
	mockInstanceWeaponRepository := new(MockActivityWeaponStatsRepository)

	mockRedisService.On("GetManifestEntity", mock.Anything, "37189237").
		Return(&types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Icon: "some/icon/url/2",
				Name: "Izanagi's Burden",
			},
			EquippingBlock: types.EquippingBlock{
				AmmoType:              1,
				EquipmentSlotTypeHash: 1498876634,
			},
		}, nil)

	weapon := model.WeaponEntity{
		WeaponHash:          37189237,
		WeaponIcon:          "some/icon/url/2",
		WeaponName:          "Izanagi's Burden",
		WeaponDamageType:    "Kinetic",
		WeaponEquipmentSlot: "Primary",
	}

	mockWeaponRepository.On("AddWeapon", &tx, weapon).
		Return(&model.WeaponEntity{}, nil)

	mockInstanceWeaponRepository.On("AddInstanceWeaponStats", &tx, mock.Anything).
		Return(nil, fmt.Errorf("Something happened when adding weapon stats to DB"))

	sut := PgcrServiceImpl{
		Redis:                         mockRedisService,
		WeaponRepository:              mockWeaponRepository,
		InstanceWeaponStatsRepository: mockInstanceWeaponRepository,
	}

	// when: addWeaponInfoToTx is called
	err := sut.addWeaponInfoToTx(&tx, &ppgcr)

	// then: an error is returned when saving instance weapon stats
	if err == nil {
		t.Fatalf("Expecting error, found none")
	}
}

func TestMapOverallStats_Success_FlawlessDuo(t *testing.T) {
	// given: a processed pgcr
	ppgcr := FlawlessDuo()

	// when: mapOverallStats is called
	result := mapOverallStats(&ppgcr)

	// then: the result should have correct overall stats
	assert := assert.New(t)
	assert.Equal(len(result), 2, "There should be only two entities")
}

func PpgcrWithWeaponStats() types.ProcessedPostGameCarnageReport {
	return types.ProcessedPostGameCarnageReport{
		InstanceId: 213123,
		PlayerInformation: []types.PlayerData{
			{
				MembershipId:          4611686018440744095,
				MembershipType:        1,
				DisplayName:           "Ceriumz",
				GlobalDisplayName:     "Ceriumz",
				GlobalDisplayNameCode: 1527,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:       2305843009263810795,
						CharacterEmblem:   2847579025,
						CharacterClass:    "Hunter",
						Kills:             376,
						Assists:           94,
						Deaths:            15,
						Kda:               28.20,
						Kdr:               25.07,
						ActivityCompleted: false,
						TimePlayedSeconds: 6953,
						AbilityInformation: types.CharacterAbilityInformation{
							GrenadeKills: 0,
							MeleeKills:   136,
							SuperKills:   3,
						},
						WeaponInformation: []types.CharacterWeaponInformation{
							{
								WeaponHash:     37189237,
								Kills:          129,
								PrecisionKills: 73,
								PrecisionRatio: 15.0,
							},
							{
								WeaponHash:     8978912,
								Kills:          31,
								PrecisionKills: 4,
								PrecisionRatio: 5.0,
							},
						},
					},
				},
			},
		},
	}
}

func FlawlessDuo() types.ProcessedPostGameCarnageReport {
	return types.ProcessedPostGameCarnageReport{
		InstanceId: 13951743181,
		PlayerInformation: []types.PlayerData{
			{
				MembershipId:          4611686018440744095,
				MembershipType:        1,
				DisplayName:           "Ceriumz",
				GlobalDisplayName:     "Ceriumz",
				GlobalDisplayNameCode: 1527,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:       2305843009263810797,
						CharacterEmblem:   2847579025,
						CharacterClass:    "Hunter",
						Kills:             376,
						Assists:           94,
						Deaths:            0,
						Kda:               28.20,
						Kdr:               25.07,
						ActivityCompleted: true,
						TimePlayedSeconds: 6953,
						AbilityInformation: types.CharacterAbilityInformation{
							GrenadeKills: 0,
							MeleeKills:   136,
							SuperKills:   3,
						},
					},
				},
			},
			{
				MembershipId:          4611686018428436415,
				MembershipType:        2,
				DisplayName:           "Dragonking890",
				GlobalDisplayName:     "Dragonking890",
				GlobalDisplayNameCode: 7731,
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:       2305843009261607826,
						CharacterEmblem:   54004491,
						CharacterClass:    "Warlock",
						Kills:             323,
						Assists:           93,
						Deaths:            0,
						Kda:               369.5,
						Kdr:               323.0,
						ActivityCompleted: true,
						TimePlayedSeconds: 1910,
						AbilityInformation: types.CharacterAbilityInformation{
							GrenadeKills: 9,
							MeleeKills:   9,
							SuperKills:   7,
						},
					},
				},
			},
		},
	}
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
						CharacterClass:  "Hunter",
					},
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "Warlock",
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
						CharacterClass:  "Titan",
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
				MembershipId:   4611686018440744095,
				MembershipType: 1,
				DisplayName:    "Ceriumz",
				PlayerCharacterInformation: []types.PlayerCharacterInformation{
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "Hunter",
					},
					{
						CharacterId:     2305843009263810795,
						CharacterEmblem: 2847579025,
						CharacterClass:  "Warlock",
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
						CharacterClass:  "Titan",
					},
				},
			},
		},
	}
}
