package processor

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"rivenbot/internal/dto"
	"rivenbot/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisService struct {
	mock.Mock
}

func (m *MockRedisService) GetManifestEntity(ctx context.Context, key string) (*dto.ManifestObject, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*dto.ManifestObject), args.Error(1)
}

// Testing solo flawless
func TestSoloFlawlessPgcr(t *testing.T) {
	// Given: a raw PGCR
	file := "../../testdata/solo_flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}

	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	response := &dto.ManifestObject{
		DisplayProperties: dto.DisplayProperties{
			Name: "Last Wish",
		},
	}
	mockedRedis.On("GetManifestEntity", mock.Anything, "2381413764").Return(response, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)

	// General lowman flags
	assert.Equal(processed.Flawless, true, "Flawless should be false")
	assert.Equal(processed.Solo, true, "Solo should be true")
	assert.Equal(processed.Duo, false, "Duo should be false")
	assert.Equal(processed.Trio, false, "Trio should be false")

	// General PGCR info
	assertPgcrFields(*processed, *pgcr, *response, assert)

	// Player information
	assertPlayers(processed.PlayerInformation, pgcr.Entries, assert)
}

// Testing duo flawless pgcr processing
func TestDuoPgcr(t *testing.T) {
	// Given: a duo flawless raw PGCR
	file := "../../testdata/duo_flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}
	slices.SortFunc(pgcr.Entries, func(a, b dto.PostGameCarnageReportEntry) int {
		return strings.Compare(a.Player.DestinyUserInfo.DisplayName, b.Player.DestinyUserInfo.DisplayName)
	})

	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	response := &dto.ManifestObject{
		DisplayProperties: dto.DisplayProperties{
			Name: "Vault of Glass: Normal",
		},
	}
	mockedRedis.On("GetManifestEntity", mock.Anything, "3881495763").Return(response, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)
	slices.SortFunc(processed.PlayerInformation, func(a, b model.PlayerInformation) int {
		return strings.Compare(a.DisplayName, b.DisplayName)
	})

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)

	assert.Equal(len(processed.PlayerInformation), 2, "There should only be 2 players")
	assert.Equal(processed.Flawless, true, "Flawless should be true")
	assert.Equal(processed.Solo, false, "Solo should be false")
	assert.Equal(processed.Duo, true, "Duo should be true")
	assert.Equal(processed.Trio, false, "Trio should be false")

	assertPgcrFields(*processed, *pgcr, *response, assert)
	assertPlayers(processed.PlayerInformation, pgcr.Entries, assert)
}

// Testing trio flawless
func TestTrioPgcr(t *testing.T) {
	// Given: a raw PGCR
	file := "../../testdata/trio_flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}

	slices.SortFunc(pgcr.Entries, func(a, b dto.PostGameCarnageReportEntry) int {
		return strings.Compare(a.CharacterId, b.CharacterId)
	})
	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	response := &dto.ManifestObject{
		DisplayProperties: dto.DisplayProperties{
			Name: "King's Fall: Normal",
		},
	}
	mockedRedis.On("GetManifestEntity", mock.Anything, "1374392663").Return(response, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)

	slices.SortFunc(processed.PlayerInformation, func(a, b model.PlayerInformation) int {
		return compareInt(a.PlayerCharacterInformation[0].CharacterId, b.PlayerCharacterInformation[0].CharacterId)
	})

	assert.Equal(len(processed.PlayerInformation), 3, "There should only be 2 players")
	assert.Equal(processed.Flawless, true, "Flawless should be true")
	assert.Equal(processed.Solo, false, "Solo should be false")
	assert.Equal(processed.Duo, false, "Duo should be false")
	assert.Equal(processed.Trio, true, "Trio should be true")

	assertPgcrFields(*processed, *pgcr, *response, assert)
	assertPlayers(processed.PlayerInformation, pgcr.Entries, assert)
}

func TestFlawlessPgcr(t *testing.T) {
	// Given: a raw PGCR with 6-man flawless
	file := "../../testdata/flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}

	slices.SortFunc(pgcr.Entries, func(a, b dto.PostGameCarnageReportEntry) int {
		return strings.Compare(a.CharacterId, b.CharacterId)
	})
	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	response := &dto.ManifestObject{
		DisplayProperties: dto.DisplayProperties{
			Name: "Crota's End: Normal",
		},
	}
	mockedRedis.On("GetManifestEntity", mock.Anything, "4179289725").Return(response, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)
	slices.SortFunc(processed.PlayerInformation, func(a, b model.PlayerInformation) int {
		return compareInt(a.PlayerCharacterInformation[0].CharacterId, b.PlayerCharacterInformation[0].CharacterId)
	})

	assert.Equal(processed.Flawless, true, "Solo should be true")
	assert.Equal(processed.Solo, false, "Solo should be false")
	assert.Equal(processed.Duo, false, "Duo should be false")
	assert.Equal(processed.Trio, false, "Trio should be true")
	assert.Equal(len(processed.PlayerInformation), 6, "There should only be 6 players")

	assertPgcrFields(*processed, *pgcr, *response, assert)
	assertPlayers(processed.PlayerInformation, pgcr.Entries, assert)
}

func getPgcr(filePath string) (*dto.PostGameCarnageReport, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pgcr dto.PostGameCarnageReportResponse

	err = json.Unmarshal(bytes, &pgcr)
	if err != nil {
		return nil, err
	}

	return &pgcr.Response, nil
}

// Utility to make assertions on general PGCRs fields
func assertPgcrFields(processed model.ProcessedPostGameCarnageReport, pgcr dto.PostGameCarnageReport, manifestObject dto.ManifestObject, assert *assert.Assertions) {
	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		assert.Errorf(err, "Error parsing time: %v", err)
	}
	instanceId, err := strconv.ParseInt(pgcr.ActivityDetails.InstanceId, 10, 64)
	if err != nil {
		assert.Error(err, "Error converting instance ID from string to int64")
	}
	endTime := startTime.Add(time.Second * time.Duration(int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)))
	raidName, raidDifficulty, err := model.Raid(manifestObject.DisplayProperties.Name)
	if err != nil {
		assert.Errorf(err, "Something bad happened when parsing the manifest activity hash")
	}

	// General PGCR info
	assert.Equal(startTime, processed.StartTime, "Start time should be the same as original pgcr.Period")
	assert.Equal(endTime, processed.EndTime, "End time should be the same when calculated")
	assert.Equal(pgcr.ActivityWasStartedFromBeginning, processed.FromBeginning, "Both should have the same beginning")
	assert.Equal(instanceId, processed.InstanceId, "Instance IDs should be the same")
	assert.Equal(raidName, processed.RaidName, "Raid name should be Vault of Glass")
	assert.Equal(raidDifficulty, processed.RaidDifficulty, "Raid difficulty should be Normal")
	assert.Equal(pgcr.ActivityDetails.ActivityHash, processed.ActivityHash, "Activity hashes should be the same")
}

// Assert players' fields in a PGCR, to be able to match players its necessary that the slices are sorted
func assertPlayers(playerInfo []model.PlayerInformation, pgcrEntries []dto.PostGameCarnageReportEntry, assert *assert.Assertions) {
	grouped, err := GroupCharacters(pgcrEntries)
	if err != nil {
		assert.Errorf(err, "Error while grouping characters for assertions")
	}
	for i, player := range playerInfo {
		membershipId, err := strconv.ParseInt(pgcrEntries[i].Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			assert.Error(err, "Something went wrong when parsing membership ID to Int64")
		}

		assert.Equal(player.DisplayName, pgcrEntries[i].Player.DestinyUserInfo.DisplayName, "Display names should be equal")
		assert.Equal(player.MembershipType, pgcrEntries[i].Player.DestinyUserInfo.MembershipType, "MembershipTypes should be equal")
		assert.Equal(player.MembershipId, membershipId, "MembershipIds should be equal")
		assert.Equal(player.GlobalDisplayName, pgcrEntries[i].Player.DestinyUserInfo.BungieGlobalDisplayName, "Global display name should be equal")
		assert.Equal(player.GlobalDisplayNameCode, pgcrEntries[i].Player.DestinyUserInfo.BungieGlobalDisplayNameCode, "Global display name code should be equal")

		// Ability information
		assert.Equal(int(pgcrEntries[i].Extended.Abilities["weaponKillsGrenade"].Basic.Value), player.PlayerCharacterInformation[0].AbilityInformation.GrenadeKills, "Grenade kills should be equal to each other")
		assert.Equal(int(pgcrEntries[i].Extended.Abilities["weaponKillsMelee"].Basic.Value), player.PlayerCharacterInformation[0].AbilityInformation.MeleeKills, "Melee kills should be equal to each other")
		assert.Equal(int(pgcrEntries[i].Extended.Abilities["weaponKillsSuper"].Basic.Value), player.PlayerCharacterInformation[0].AbilityInformation.SuperKills, "Super kills should be equal to each other")

		// Character stuff
		pgcrCharacters := grouped[player.MembershipId]
		slices.SortFunc(pgcrCharacters, func(a, b dto.PostGameCarnageReportEntry) int {
			ai, err := strconv.ParseInt(a.CharacterId, 10, 64)
			if err != nil {
				assert.Error(err, "Error parsing A characterId for sorting")
			}
			bi, err := strconv.ParseInt(b.CharacterId, 10, 64)
			if err != nil {
				assert.Error(err, "Error parsing B characterId for sorting")
			}
			return compareInt(ai, bi)
		})

		slices.SortFunc(player.PlayerCharacterInformation, func(a, b model.PlayerCharacterInformation) int {
			return compareInt(a.CharacterId, b.CharacterId)
		})
		assertPlayerCharacters(player.PlayerCharacterInformation, grouped[player.MembershipId], assert)
	}
}

// Assert all the data in player characters. Arrays must be sorted to work effectively
func assertPlayerCharacters(processed []model.PlayerCharacterInformation, pgcr []dto.PostGameCarnageReportEntry, assert *assert.Assertions) {
	for i, playerCharacter := range processed {
		characterId, err := strconv.ParseInt(pgcr[i].CharacterId, 10, 64)
		if err != nil {
			assert.Error(err, "Something went wrong when parsing character ID to Int64")
		}
		characterClass := model.CharacterClass(pgcr[i].Player.CharacterClass)
		assert.Equal(int(pgcr[i].Values["kills"].Basic.Value), playerCharacter.Kills, "Kills should match")
		assert.Equal(int(pgcr[i].Values["deaths"].Basic.Value), playerCharacter.Deaths, "Deaths should match")
		assert.Equal(int(pgcr[i].Values["assists"].Basic.Value), playerCharacter.Assists, "Assists should match")
		assert.Equal(pgcr[i].Values["killsDeathsAssists"].Basic.Value, playerCharacter.Kda, "KDA should match")
		assert.Equal(pgcr[i].Values["killsDeathsRatio"].Basic.Value, playerCharacter.Kdr, "KDR should match")
		assert.Equal(int(pgcr[i].Values["timePlayedSeconds"].Basic.Value), playerCharacter.TimePlayedSeconds, "TimePlayedSeconds should match")
		assert.Equal(characterId, playerCharacter.CharacterId, "CharacterID should match between result and original")
		assert.Equal(characterClass, playerCharacter.CharacterClass, "Character class should match")
		assert.Equal(pgcr[i].Player.EmblemHash, playerCharacter.CharacterEmblem, "Emblem hashes should match")

		slices.SortFunc(playerCharacter.WeaponInformation, func(a, b model.CharacterWeaponInformation) int {
			return compareInt(a.WeaponHash, b.WeaponHash)
		})

		slices.SortFunc(pgcr[i].Extended.Weapons, func(a, b dto.WeaponInformation) int {
			return compareInt(a.ReferenceId, b.ReferenceId)
		})

		assertPlayerWeapons(playerCharacter.WeaponInformation, pgcr[i].Extended.Weapons, assert)
	}
}

// Assert all weapons for a player character. Arrays must be sorted to work effectively
func assertPlayerWeapons(processedWeapons []model.CharacterWeaponInformation, pgcrWeapons []dto.WeaponInformation, assert *assert.Assertions) {
	for i, weapon := range processedWeapons {
		assert.Equal(weapon.WeaponHash, pgcrWeapons[i].ReferenceId, "Weapon hashes should match")
		assert.Equal(weapon.Kills, int(pgcrWeapons[i].Values["uniqueWeaponKills"].Basic.Value))
		assert.Equal(weapon.PrecisionKills, int(pgcrWeapons[i].Values["uniqueWeaponPrecisionKills"].Basic.Value))
		assert.Equal(weapon.PrecisionRatio, pgcrWeapons[i].Values["uniqueWeaponKillsPrecisionKills"].Basic.Value)
	}
}

func compareInt(a, b int64) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}
