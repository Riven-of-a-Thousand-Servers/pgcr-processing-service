package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"

	types "github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	enums "github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisService struct {
	mock.Mock
}

func (m *MockRedisService) GetManifestEntity(ctx context.Context, key string) (*types.ManifestObject, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*types.ManifestObject), args.Error(1)
}

var pgcrTests = map[string]struct {
	inputFile string
	response  *types.ManifestObject
	flawless  bool
	solo      bool
	duo       bool
	trio      bool
	size      int
}{
	"solo_flawless_pgcr": {
		inputFile: "solo_flawless_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Last Wish",
			},
		},
		size:     1,
		flawless: true,
		solo:     true,
		duo:      false,
		trio:     false,
	},
	"duo_flawless_pgcr": {
		inputFile: "duo_flawless_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Vault of Glass: Normal",
			},
		},
		size:     2,
		flawless: true,
		solo:     false,
		duo:      true,
		trio:     false,
	},
	"trio_flawless_pgcr": {
		inputFile: "trio_flawless_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "King's Fall: Normal",
			},
		},
		size:     3,
		flawless: true,
		solo:     false,
		duo:      false,
		trio:     true,
	},
	"flawless_pgcr": {
		inputFile: "flawless_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Crota's End: Normal",
			},
		},
		size:     6,
		flawless: true,
		solo:     false,
		duo:      false,
		trio:     false,
	},
	"solo_pgcr": {
		inputFile: "solo_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Root of Nightmares: Standard",
			},
		},
		size:     1,
		flawless: false,
		solo:     true,
		duo:      false,
		trio:     false,
	},
	"duo_pgcr": {
		inputFile: "duo_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Garden of Salvation",
			},
		},
		size:     2,
		flawless: false,
		solo:     false,
		duo:      true,
		trio:     false,
	},
	"trio_pgcr": {
		inputFile: "trio_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Root of Nightmares: Standard",
			},
		},
		size:     3,
		flawless: false,
		solo:     false,
		duo:      false,
		trio:     true,
	},
	"uncomplete_pgcr": {
		inputFile: "not_completed_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Root of Nightmares: Standard",
			},
		},
		size:     6,
		flawless: false,
		solo:     false,
		duo:      false,
		trio:     false,
	},
	"various_characters_on_player_pgcr": {
		inputFile: "various_character_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Root of Nightmares: Standard",
			},
		},
		size:     4,
		flawless: false,
		solo:     false,
		duo:      false,
		trio:     false,
	},
	"beyond_light_pgcr": {
		inputFile: "beyond_light_pgcr.json",
		response: &types.ManifestObject{
			DisplayProperties: types.DisplayProperties{
				Name: "Deep Stone Crypt",
			},
		},
		size:     10,
		flawless: false,
		solo:     false,
		duo:      false,
		trio:     false,
	},
}

func TestPgcrMapping(t *testing.T) {
	for test, params := range pgcrTests {
		t.Run(test, func(t *testing.T) {
			pgcr, err := getPgcr(params.inputFile)
			if err != nil {
				t.Errorf("Failed to get file [%s]. %v", params.inputFile, err)
			}

			mockedRedis := new(MockRedisService)

			// Mock manifest calls
			activityId := pgcr.ActivityDetails.ActivityHash
			mockedRedis.On("GetManifestEntity", mock.Anything, strconv.Itoa(int(activityId))).Return(params.response, nil)

			processor := PgcrMapper{
				RedisClient: mockedRedis,
			}

			_, processed, err := processor.Map(pgcr)

			assert := assert.New(t)
			slices.SortFunc(processed.PlayerInformation, func(a, b types.PlayerData) int {
				return compareInt(a.MembershipId, b.MembershipId)
			})

			// General lowman flags
			assert.Equal(processed.Flawless, params.flawless, fmt.Sprintf("Flawless should be %v", params.flawless))
			assert.Equal(processed.Solo, params.solo, fmt.Sprintf("Solo should be %v", params.solo))
			assert.Equal(processed.Duo, params.duo, fmt.Sprintf("Duo should be %v", params.duo))
			assert.Equal(processed.Trio, params.trio, fmt.Sprintf("Trio should be %v", params.trio))

			// Player Information size is correct
			assert.Equal(len(processed.PlayerInformation), params.size, fmt.Sprintf("Player Information size should be %d", params.size))

			// Assert PGCR info
			assertPgcrFields(*processed, *pgcr, *params.response, assert)

			// Assert player information
			assertPlayers(processed.PlayerInformation, pgcr.Entries, assert)
		})
	}
}

var freshness = map[string]struct {
	inputFile string
	fresh     bool
}{
	"beyong_light_pgcr": {
		inputFile: "beyond_light_pgcr.json",
		fresh:     false,
	},
	"witch_queen_pgcr": {
		inputFile: "witch_queen_pgcr.json",
		fresh:     false,
	},
	"pre_beyond_light_non_fresh_pgcr": {
		inputFile: "pre_beyond_light_non_fresh_pgcr.json",
		fresh:     false,
	},
	"pre_beyond_light_fresh_pgcr": {
		inputFile: "pre_beyond_light_fresh_pgcr.json",
		fresh:     true,
	},
}

// Test various pgcr to see if they should be fresh or not
func TestPgcrFreshness(t *testing.T) {
	for test, params := range freshness {
		t.Run(test, func(t *testing.T) {
			pgcr, err := getPgcr(params.inputFile)
			if err != nil {
				t.Errorf("Failed to get file [%s]. %v", params.inputFile, err)
			}

			mockedRedis := new(MockRedisService)

			// Mock manifest calls
			activityId := pgcr.ActivityDetails.ActivityHash
			response := &types.ManifestObject{
				DisplayProperties: types.DisplayProperties{
					Name: "Last Wish",
				},
			}
			mockedRedis.On("GetManifestEntity", mock.Anything, strconv.Itoa(int(activityId))).Return(response, nil)

			processor := PgcrMapper{
				RedisClient: mockedRedis,
			}

			_, processed, err := processor.Map(pgcr)

			assert := assert.New(t)

			assert.Equal(processed.FromBeginning, params.fresh, fmt.Sprintf("FromBeginning attribute should be %v", params.fresh))
		})
	}
}

// Utility to retrieve a pgcr json as test data
func getPgcr(file string) (*types.PostGameCarnageReport, error) {
	bytes, err := os.ReadFile("../testdata/" + file)
	if err != nil {
		return nil, err
	}

	var pgcr types.PostGameCarnageReportResponse

	err = json.Unmarshal(bytes, &pgcr)
	if err != nil {
		return nil, err
	}

	return &pgcr.Response, nil
}

// Utility to make assertions on general PGCRs fields
func assertPgcrFields(processed types.ProcessedPostGameCarnageReport, pgcr types.PostGameCarnageReport, manifestObject types.ManifestObject, assert *assert.Assertions) {
	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		assert.Errorf(err, "Error parsing time: %v", err)
	}
	instanceId, err := strconv.ParseInt(pgcr.ActivityDetails.InstanceId, 10, 64)
	if err != nil {
		assert.Error(err, "Error converting instance ID from string to int64")
	}
	endTime := startTime.Add(time.Second * time.Duration(int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)))
	raidName, raidDifficulty, err := enums.GetRaidAndDifficulty(manifestObject.DisplayProperties.Name)
	if err != nil {
		assert.Errorf(err, "Something bad happened when parsing the manifest activity hash")
	}

	// General PGCR info
	assert.Equal(startTime, processed.StartTime, "Start time should be the same as the original pgcr.Period")
	assert.Equal(endTime, processed.EndTime, "End time should be the same when calculated")
	assert.Equal(pgcr.ActivityWasStartedFromBeginning, processed.FromBeginning, "Both should have the same beginning")
	assert.Equal(instanceId, processed.InstanceId, "Instance IDs should be the same")
	assert.Equal(raidName, processed.RaidName, fmt.Sprintf("Raid name should be %s", raidName))
	assert.Equal(raidDifficulty, processed.RaidDifficulty, fmt.Sprintf("Raid difficulty should be %s", raidDifficulty))
	assert.Equal(pgcr.ActivityDetails.ActivityHash, processed.ActivityHash, "Activity hashes should be the same")
}

// Assert players' fields in a PGCR, to be able to match players its necessary that the slices are sorted
func assertPlayers(playerData []types.PlayerData, pgcrEntries []types.PostGameCarnageReportEntry, assert *assert.Assertions) {
	grouped, err := GroupCharacters(pgcrEntries)
	if err != nil {
		assert.Errorf(err, "Error while grouping characters for assertions")
	}

	// Sort all the keys, which are the membership Ids of all players in the activity
	keys := make([]int64, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(a, b int64) int {
		return compareInt(a, b)
	})

	for i, player := range playerData {
		membershipId, err := strconv.ParseInt(grouped[keys[i]][0].Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			assert.Error(err, "Something went wrong when parsing membership ID to Int64")
		}

		assert.Equal(player.DisplayName, grouped[keys[i]][0].Player.DestinyUserInfo.DisplayName, "Display names should be equal")
		assert.Equal(player.MembershipType, grouped[keys[i]][0].Player.DestinyUserInfo.MembershipType, "MembershipTypes should be equal")
		assert.Equal(player.MembershipId, membershipId, "MembershipIds should be equal")
		assert.Equal(player.GlobalDisplayName, grouped[keys[i]][0].Player.DestinyUserInfo.BungieGlobalDisplayName, "Global display name should be equal")
		assert.Equal(player.GlobalDisplayNameCode, grouped[keys[i]][0].Player.DestinyUserInfo.BungieGlobalDisplayNameCode, "Global display name code should be equal")

		// Character stuff
		pgcrCharacters := grouped[player.MembershipId]
		slices.SortFunc(pgcrCharacters, func(a, b types.PostGameCarnageReportEntry) int {
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

		slices.SortFunc(player.PlayerCharacterInformation, func(a, b types.PlayerCharacterInformation) int {
			return compareInt(a.CharacterId, b.CharacterId)
		})
		assertPlayerCharacters(player.PlayerCharacterInformation, grouped[player.MembershipId], assert)
	}
}

// Assert all the data in player characters. Arrays must be sorted to work effectively
func assertPlayerCharacters(processed []types.PlayerCharacterInformation, pgcr []types.PostGameCarnageReportEntry, assert *assert.Assertions) {
	for i, playerCharacter := range processed {
		characterId, err := strconv.ParseInt(pgcr[i].CharacterId, 10, 64)
		if err != nil {
			assert.Error(err, "Something went wrong when parsing character ID to Int64")
		}
		characterClass := types.CharacterClass(pgcr[i].Player.CharacterClass)
		assert.Equal(len(processed), len(pgcr), "Both slices should contain the same amount of characters for a player")
		assert.Equal(int(pgcr[i].Values["kills"].Basic.Value), playerCharacter.Kills, "Kills should match")
		assert.Equal(pgcr[i].Values["completed"].Basic.Value == 1.0, playerCharacter.ActivityCompleted, "Completed should be correct for the player")
		assert.Equal(int(pgcr[i].Values["deaths"].Basic.Value), playerCharacter.Deaths, "Deaths should match")
		assert.Equal(int(pgcr[i].Values["assists"].Basic.Value), playerCharacter.Assists, "Assists should match")
		assert.Equal(pgcr[i].Values["killsDeathsAssists"].Basic.Value, playerCharacter.Kda, "KDA should match")
		assert.Equal(pgcr[i].Values["killsDeathsRatio"].Basic.Value, playerCharacter.Kdr, "KDR should match")
		assert.Equal(int(pgcr[i].Values["timePlayedSeconds"].Basic.Value), playerCharacter.TimePlayedSeconds, "TimePlayedSeconds should match")
		assert.Equal(characterId, playerCharacter.CharacterId, "CharacterID should match between result and original")
		assert.Equal(characterClass, playerCharacter.CharacterClass, "Character class should match")
		assert.Equal(pgcr[i].Player.EmblemHash, playerCharacter.CharacterEmblem, "Emblem hashes should match")

		// Ability information
		assert.Equal(int(pgcr[i].Extended.Abilities["weaponKillsGrenade"].Basic.Value), playerCharacter.AbilityInformation.GrenadeKills, "Grenade kills should be equal to each other")
		assert.Equal(int(pgcr[i].Extended.Abilities["weaponKillsMelee"].Basic.Value), playerCharacter.AbilityInformation.MeleeKills, "Melee kills should be equal to each other")
		assert.Equal(int(pgcr[i].Extended.Abilities["weaponKillsSuper"].Basic.Value), playerCharacter.AbilityInformation.SuperKills, "Super kills should be equal to each other")

		slices.SortFunc(playerCharacter.WeaponInformation, func(a, b types.CharacterWeaponInformation) int {
			return compareInt(a.WeaponHash, b.WeaponHash)
		})

		slices.SortFunc(pgcr[i].Extended.Weapons, func(a, b types.WeaponInformation) int {
			return compareInt(a.ReferenceId, b.ReferenceId)
		})

		assertPlayerWeapons(playerCharacter.WeaponInformation, pgcr[i].Extended.Weapons, assert)
	}
}

// Assert all weapons for a player character. Arrays must be sorted to work effectively
func assertPlayerWeapons(processedWeapons []types.CharacterWeaponInformation, pgcrWeapons []types.WeaponInformation, assert *assert.Assertions) {
	for i, weapon := range processedWeapons {
		assert.Equal(weapon.WeaponHash, pgcrWeapons[i].ReferenceId, "Weapon hashes should match")
		assert.Equal(weapon.Kills, int(pgcrWeapons[i].Values["uniqueWeaponKills"].Basic.Value), "Weapon kills should match")
		assert.Equal(weapon.PrecisionKills, int(pgcrWeapons[i].Values["uniqueWeaponPrecisionKills"].Basic.Value), "Weapon precision kills should match")
		assert.Equal(weapon.PrecisionRatio, pgcrWeapons[i].Values["uniqueWeaponKillsPrecisionKills"].Basic.Value, "Weapon precision ratio should match")
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
