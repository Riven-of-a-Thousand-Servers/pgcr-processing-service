package pgcr

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
	file := "../testdata/solo_flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}

	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	mockedRedis.On("GetManifestEntity", mock.Anything, "2381413764").Return(
		&dto.ManifestObject{
			DisplayProperties: dto.DisplayProperties{
				Name: "Last Wish",
			},
		}, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)

	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		t.Errorf("Error parsing time: %v", err)
	}

	endTime := startTime.Add(time.Second * time.Duration(int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)))

	// General PGCR info
	assert.Equal(processed.StartTime, startTime, "Start time should be the same as original pgcr.Period")
	assert.Equal(processed.EndTime, endTime, "End time should be the same when calculated")
	assert.Equal(processed.FromBeginning, pgcr.ActivityWasStartedFromBeginning, "Both should have the same beginning")
	assert.Equal(processed.InstanceId, int64(14287236297), "Instance IDs should be the same")
	assert.Equal(processed.RaidName, model.LAST_WISH, "Raid name should be valid and equal")
	assert.Equal(processed.RaidDifficulty, model.NORMAL, "Raid difficulty should be valid and equal")
	assert.Equal(processed.ActivityHash, pgcr.ActivityDetails.ActivityHash, "Activity hashes should be the same")
	assert.Equal(processed.Flawless, true, "Flawless should be false")
	assert.Equal(processed.Solo, true, "Solo should be true")
	assert.Equal(processed.Duo, false, "Duo should be false")
	assert.Equal(processed.Trio, false, "Trio should be false")

	// Player information
	membershipId, err := strconv.ParseInt(pgcr.Entries[0].Player.DestinyUserInfo.MembershipId, 10, 64)
	if err != nil {
		t.Error("Something went wrong when parsing membership ID to Int64")
	}
	assert.Equal(len(processed.PlayerInformation), 1, "Length of player info should be 1")
	assert.Equal(processed.PlayerInformation[0].DisplayName, pgcr.Entries[0].Player.DestinyUserInfo.DisplayName, "Display names should be equal")
	assert.Equal(processed.PlayerInformation[0].MembershipType, pgcr.Entries[0].Player.DestinyUserInfo.MembershipType, "MembershipTypes should be equal")
	assert.Equal(processed.PlayerInformation[0].MembershipId, membershipId, "MembershipIds should be equal")
	assert.Equal(processed.PlayerInformation[0].GlobalDisplayName, pgcr.Entries[0].Player.DestinyUserInfo.BungieGlobalDisplayName, "Global display name should be equal")
	assert.Equal(processed.PlayerInformation[0].GlobalDisplayNameCode, pgcr.Entries[0].Player.DestinyUserInfo.BungieGlobalDisplayNameCode, "Global display name code should be equal")

	// Character stuff
	characterClass := model.CharacterClass(pgcr.Entries[0].Player.CharacterClass)
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].Kills, int(pgcr.Entries[0].Values["kills"].Basic.Value), "Kills should match")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].Deaths, int(pgcr.Entries[0].Values["deaths"].Basic.Value), "Deaths should match")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].Assists, int(pgcr.Entries[0].Values["assists"].Basic.Value), "Assists should match")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].Kda, pgcr.Entries[0].Values["killsDeathsAssists"].Basic.Value, "KDA should match")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].Kdr, pgcr.Entries[0].Values["killsDeathsRatio"].Basic.Value, "KDR should match")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].TimePlayedSeconds, int(pgcr.Entries[0].Values["timePlayedSeconds"].Basic.Value), "TimePlayedSeconds should match")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].CharacterId, int64(2305843009468984093), "CharacterID should match between result and original")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].CharacterClass, characterClass, "CharacterEmblem")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].CharacterEmblem, pgcr.Entries[0].Player.EmblemHash, "CharacterID should match between result and original")

	// Ability & weapon info
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].WeaponInformation[0].WeaponHash, pgcr.Entries[0].Extended.Weapons[0].ReferenceId, "Weapon should have matching weaponID")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].WeaponInformation[0].Kills, int(pgcr.Entries[0].Extended.Weapons[0].Values["uniqueWeaponKills"].Basic.Value), "Weapon should have matching kills")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].WeaponInformation[0].PrecisionKills, int(pgcr.Entries[0].Extended.Weapons[0].Values["uniqueWeaponPrecisionKills"].Basic.Value), "Weapon should have matching precision kills")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].WeaponInformation[0].PrecisionRatio, pgcr.Entries[0].Extended.Weapons[0].Values["uniqueWeaponKillsPrecisionKills"].Basic.Value, "Weapon should have matching precision ratio")

	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].AbilityInformation.GrenadeKills, int(pgcr.Entries[0].Extended.Abilities["weaponKillsGrenade"].Basic.Value), "Grenade kills should be equal to each other")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].AbilityInformation.MeleeKills, int(pgcr.Entries[0].Extended.Abilities["weaponKillsMelee"].Basic.Value), "Melee kills should be equal to each other")
	assert.Equal(processed.PlayerInformation[0].PlayerCharacterInformation[0].AbilityInformation.SuperKills, int(pgcr.Entries[0].Extended.Abilities["weaponKillsSuper"].Basic.Value), "Super kills should be equal to each other")
}

// Testing duo flawless pgcr processing
func TestDuoPgcr(t *testing.T) {
	// Given: a duo flawless raw PGCR
	file := "../testdata/duo_flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}

	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	mockedRedis.On("GetManifestEntity", mock.Anything, "3881495763").Return(
		&dto.ManifestObject{
			DisplayProperties: dto.DisplayProperties{
				Name: "Vault of Glass: Normal",
			},
		}, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)

	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		t.Errorf("Error parsing time: %v", err)
	}

	endTime := startTime.Add(time.Second * time.Duration(int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)))

	// General PGCR info
	assert.Equal(processed.StartTime, startTime, "Start time should be the same as original pgcr.Period")
	assert.Equal(processed.EndTime, endTime, "End time should be the same when calculated")
	assert.Equal(processed.FromBeginning, pgcr.ActivityWasStartedFromBeginning, "Both should have the same beginning")
	assert.Equal(processed.InstanceId, int64(10902450421), "Instance IDs should be the same")
	assert.Equal(processed.RaidName, model.VAULT_OF_GLASS, "Raid name should be Vault of Glass")
	assert.Equal(processed.RaidDifficulty, model.NORMAL, "Raid difficulty should be Normal")
	assert.Equal(processed.ActivityHash, pgcr.ActivityDetails.ActivityHash, "Activity hashes should be the same")
	assert.Equal(processed.Flawless, true, "Flawless should be true")
	assert.Equal(processed.Solo, false, "Solo should be false")
	assert.Equal(processed.Duo, true, "Duo should be true")
	assert.Equal(processed.Trio, false, "Trio should be false")

	// Player information
	assert.Equal(len(processed.PlayerInformation), 2, "There should only be 2 players")
	for i, player := range processed.PlayerInformation {
		membershipId, err := strconv.ParseInt(pgcr.Entries[i].Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			t.Error("Something went wrong when parsing membership ID to Int64")
		}
		characterId, err := strconv.ParseInt(pgcr.Entries[i].CharacterId, 10, 64)
		if err != nil {
			t.Error("Something went wrong when parsing character ID to Int64")
		}
		assert.Equal(player.DisplayName, pgcr.Entries[i].Player.DestinyUserInfo.DisplayName, "Display names should be equal")
		assert.Equal(player.MembershipType, pgcr.Entries[i].Player.DestinyUserInfo.MembershipType, "MembershipTypes should be equal")
		assert.Equal(player.MembershipId, membershipId, "MembershipIds should be equal")
		assert.Equal(player.GlobalDisplayName, pgcr.Entries[i].Player.DestinyUserInfo.BungieGlobalDisplayName, "Global display name should be equal")
		assert.Equal(player.GlobalDisplayNameCode, pgcr.Entries[i].Player.DestinyUserInfo.BungieGlobalDisplayNameCode, "Global display name code should be equal")

		// Character stuff
		characterClass := model.CharacterClass(pgcr.Entries[i].Player.CharacterClass)
		assert.Equal(player.PlayerCharacterInformation[0].Kills, int(pgcr.Entries[i].Values["kills"].Basic.Value), "Kills should match")
		assert.Equal(player.PlayerCharacterInformation[0].Deaths, int(pgcr.Entries[i].Values["deaths"].Basic.Value), "Deaths should match")
		assert.Equal(player.PlayerCharacterInformation[0].Assists, int(pgcr.Entries[i].Values["assists"].Basic.Value), "Assists should match")
		assert.Equal(player.PlayerCharacterInformation[0].Kda, pgcr.Entries[i].Values["killsDeathsAssists"].Basic.Value, "KDA should match")
		assert.Equal(player.PlayerCharacterInformation[0].Kdr, pgcr.Entries[i].Values["killsDeathsRatio"].Basic.Value, "KDR should match")
		assert.Equal(player.PlayerCharacterInformation[0].TimePlayedSeconds, int(pgcr.Entries[i].Values["timePlayedSeconds"].Basic.Value), "TimePlayedSeconds should match")
		assert.Equal(player.PlayerCharacterInformation[0].CharacterId, characterId, "CharacterID should match between result and original")
		assert.Equal(player.PlayerCharacterInformation[0].CharacterClass, characterClass, "CharacterEmblem")
		assert.Equal(player.PlayerCharacterInformation[0].CharacterEmblem, pgcr.Entries[i].Player.EmblemHash, "CharacterID should match between result and original")

		// Ability & weapon info
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].WeaponHash, pgcr.Entries[i].Extended.Weapons[0].ReferenceId, "Weapon should have matching weaponID")
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].Kills, int(pgcr.Entries[i].Extended.Weapons[0].Values["uniqueWeaponKills"].Basic.Value), "Weapon should have matching kills")
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].PrecisionKills, int(pgcr.Entries[i].Extended.Weapons[0].Values["uniqueWeaponPrecisionKills"].Basic.Value), "Weapon should have matching precision kills")
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].PrecisionRatio, pgcr.Entries[i].Extended.Weapons[0].Values["uniqueWeaponKillsPrecisionKills"].Basic.Value, "Weapon should have matching precision ratio")

		assert.Equal(player.PlayerCharacterInformation[0].AbilityInformation.GrenadeKills, int(pgcr.Entries[i].Extended.Abilities["weaponKillsGrenade"].Basic.Value), "Grenade kills should be equal to each other")
		assert.Equal(player.PlayerCharacterInformation[0].AbilityInformation.MeleeKills, int(pgcr.Entries[i].Extended.Abilities["weaponKillsMelee"].Basic.Value), "Melee kills should be equal to each other")
		assert.Equal(player.PlayerCharacterInformation[0].AbilityInformation.SuperKills, int(pgcr.Entries[i].Extended.Abilities["weaponKillsSuper"].Basic.Value), "Super kills should be equal to each other")
	}
}

// Testing trio flawless
func TestTrioPgcr(t *testing.T) {
	// Given: a raw PGCR
	file := "../testdata/trio_flawless_pgcr.json"
	pgcr, err := getPgcr(file)
	if err != nil {
		t.Errorf("Failed to get file [%s]. %v", file, err)
	}

	slices.SortFunc(pgcr.Entries, func(a, b dto.PostGameCarnageReportEntry) int {
		return strings.Compare(a.CharacterId, b.CharacterId)
	})
	mockedRedis := new(MockRedisService)

	// Mock manifest calls
	mockedRedis.On("GetManifestEntity", mock.Anything, "1374392663").Return(
		&dto.ManifestObject{
			DisplayProperties: dto.DisplayProperties{
				Name: "King's Fall: Normal",
			},
		}, nil)

	processor := PGCRProcessor{
		redisClient: mockedRedis,
	}

	// When: Process is called
	_, processed, err := processor.Process(pgcr)

	// Then: the return values are valid and processing went smooth
	assert := assert.New(t)

	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		t.Errorf("Error parsing time: %v", err)
	}

	slices.SortFunc(processed.PlayerInformation, func(a, b model.PlayerInformation) int {
		if a.PlayerCharacterInformation[0].CharacterId > a.PlayerCharacterInformation[0].CharacterId {
			return 1
		} else if a.PlayerCharacterInformation[0].CharacterId < a.PlayerCharacterInformation[0].CharacterId {
			return -1
		} else {
			return 0
		}
	})

	endTime := startTime.Add(time.Second * time.Duration(int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)))

	// General PGCR info
	assert.Equal(processed.StartTime, startTime, "Start time should be the same as original pgcr.Period")
	assert.Equal(processed.EndTime, endTime, "End time should be the same when calculated")
	assert.Equal(processed.FromBeginning, pgcr.ActivityWasStartedFromBeginning, "Both should have the same beginning")
	assert.Equal(processed.InstanceId, int64(11756061781), "Instance IDs should be the same")
	assert.Equal(processed.RaidName, model.KINGS_FALL, "Raid name should be Vault of Glass")
	assert.Equal(processed.RaidDifficulty, model.NORMAL, "Raid difficulty should be Normal")
	assert.Equal(processed.ActivityHash, pgcr.ActivityDetails.ActivityHash, "Activity hashes should be the same")
	assert.Equal(processed.Flawless, true, "Flawless should be true")
	assert.Equal(processed.Solo, false, "Solo should be false")
	assert.Equal(processed.Duo, false, "Duo should be false")
	assert.Equal(processed.Trio, true, "Trio should be true")

	assert.Equal(len(processed.PlayerInformation), 3, "There should only be 2 players")
	for i, player := range processed.PlayerInformation {
		membershipId, err := strconv.ParseInt(pgcr.Entries[i].Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			t.Error("Something went wrong when parsing membership ID to Int64")
		}
		characterId, err := strconv.ParseInt(pgcr.Entries[i].CharacterId, 10, 64)
		if err != nil {
			t.Error("Something went wrong when parsing character ID to Int64")
		}
		assert.Equal(player.DisplayName, pgcr.Entries[i].Player.DestinyUserInfo.DisplayName, "Display names should be equal")
		assert.Equal(player.MembershipType, pgcr.Entries[i].Player.DestinyUserInfo.MembershipType, "MembershipTypes should be equal")
		assert.Equal(player.MembershipId, membershipId, "MembershipIds should be equal")
		assert.Equal(player.GlobalDisplayName, pgcr.Entries[i].Player.DestinyUserInfo.BungieGlobalDisplayName, "Global display name should be equal")
		assert.Equal(player.GlobalDisplayNameCode, pgcr.Entries[i].Player.DestinyUserInfo.BungieGlobalDisplayNameCode, "Global display name code should be equal")

		// Character stuff
		characterClass := model.CharacterClass(pgcr.Entries[i].Player.CharacterClass)
		assert.Equal(player.PlayerCharacterInformation[0].Kills, int(pgcr.Entries[i].Values["kills"].Basic.Value), "Kills should match")
		assert.Equal(player.PlayerCharacterInformation[0].Deaths, int(pgcr.Entries[i].Values["deaths"].Basic.Value), "Deaths should match")
		assert.Equal(player.PlayerCharacterInformation[0].Assists, int(pgcr.Entries[i].Values["assists"].Basic.Value), "Assists should match")
		assert.Equal(player.PlayerCharacterInformation[0].Kda, pgcr.Entries[i].Values["killsDeathsAssists"].Basic.Value, "KDA should match")
		assert.Equal(player.PlayerCharacterInformation[0].Kdr, pgcr.Entries[i].Values["killsDeathsRatio"].Basic.Value, "KDR should match")
		assert.Equal(player.PlayerCharacterInformation[0].TimePlayedSeconds, int(pgcr.Entries[i].Values["timePlayedSeconds"].Basic.Value), "TimePlayedSeconds should match")
		assert.Equal(player.PlayerCharacterInformation[0].CharacterId, characterId, "CharacterID should match between result and original")
		assert.Equal(player.PlayerCharacterInformation[0].CharacterClass, characterClass, "CharacterEmblem")
		assert.Equal(player.PlayerCharacterInformation[0].CharacterEmblem, pgcr.Entries[i].Player.EmblemHash, "CharacterID should match between result and original")

		// Ability & weapon info
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].WeaponHash, pgcr.Entries[i].Extended.Weapons[0].ReferenceId, "Weapon should have matching weaponID")
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].Kills, int(pgcr.Entries[i].Extended.Weapons[0].Values["uniqueWeaponKills"].Basic.Value), "Weapon should have matching kills")
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].PrecisionKills, int(pgcr.Entries[i].Extended.Weapons[0].Values["uniqueWeaponPrecisionKills"].Basic.Value), "Weapon should have matching precision kills")
		assert.Equal(player.PlayerCharacterInformation[0].WeaponInformation[0].PrecisionRatio, pgcr.Entries[i].Extended.Weapons[0].Values["uniqueWeaponKillsPrecisionKills"].Basic.Value, "Weapon should have matching precision ratio")

		assert.Equal(player.PlayerCharacterInformation[0].AbilityInformation.GrenadeKills, int(pgcr.Entries[i].Extended.Abilities["weaponKillsGrenade"].Basic.Value), "Grenade kills should be equal to each other")
		assert.Equal(player.PlayerCharacterInformation[0].AbilityInformation.MeleeKills, int(pgcr.Entries[i].Extended.Abilities["weaponKillsMelee"].Basic.Value), "Melee kills should be equal to each other")
		assert.Equal(player.PlayerCharacterInformation[0].AbilityInformation.SuperKills, int(pgcr.Entries[i].Extended.Abilities["weaponKillsSuper"].Basic.Value), "Super kills should be equal to each other")
	}
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
