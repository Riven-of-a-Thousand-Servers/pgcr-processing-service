package pgcr

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

  "github.com/redis/go-redis/v9"

  utils "rivenbot/pkg/utils"
	r  "rivenbot/pkg/redis"
	dto "rivenbot/types/dto"
	entity "rivenbot/types/entity"
)

const (
	sotpHash1   int64  = 548750096
	sotpHash2   int64  = 2812525063
	pstTimezone string = "America/Los_Angeles"
)

var (
	beyondLightStart = time.Date(2020, time.November, 10, 9, 0, 0, 0, time.FixedZone("PST", -8*60*60))
	witchQueenStart  = time.Date(2022, time.February, 22, 9, 0, 0, 0, time.FixedZone("PST", -8*60*60))
	hauntedStart     = time.Date(2022, time.May, 24, 10, 0, 0, 0, time.FixedZone("PDT", -7*60*60))
)

var leviHashes = map[int64]bool{
	2693136600: true, 2693136601: true, 2693136602: true,
	2693136603: true, 2693136604: true, 2693136605: true,
	89727599: true, 287649202: true, 1699948563: true, 1875726950: true,
	3916343513: true, 4039317196: true, 417231112: true, 508802457: true,
	757116822: true, 771164842: true, 1685065161: true, 1800508819: true,
	2449714930: true, 3446541099: true, 4206123728: true, 3912437239: true,
	3879860661: true, 3857338478: true,
}

// This method processes the PGCR into a format thats more suitable
func Process(pgcr *dto.PostGameCarnageReport) (*entity.RaidPgcr, error){
	processedPgcr, err := createProcessedPgcr(pgcr)
	if err != nil {
		log.Fatal(err)
    return nil, err
	}

  compressedData, err := Compress(processedPgcr)
  if err != nil {
    return nil, err
  }
  raidPgcr := entity.RaidPgcr{
    InstanceId: processedPgcr.InstanceId,
    Timestamp: processedPgcr.StartTime,
    Blob: compressedData,
  }
  return &raidPgcr, nil
}

func createProcessedPgcr(pgcr *dto.PostGameCarnageReport) (*entity.ProcessedPostGameCarnageReport, error) {
	entity := new(entity.ProcessedPostGameCarnageReport)

  // Calculate start and end time
	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		return nil, err
	}

	if len(pgcr.Entries) == 0 {
		error := fmt.Errorf("No entries in the PostGameCarnageReport, unable to determine the end time of the activity")
		return nil, error
	}

	activityDurationSeconds := int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)
	endTime := startTime.Add(time.Second * time.Duration(activityDurationSeconds))

	*&entity.StartTime = startTime
	*&entity.EndTime = endTime
	*&entity.InstanceId = pgcr.ActivityDetails.InstanceId

  redisClient, err := r.CreateClient()
  if err != nil {
    return nil, err
  }

  maniestResponse, err := r.GetManifestEntity(redisClient,
    strconv.Itoa(int(pgcr.ActivityDetails.ActivityHash)))
  if err != nil {
    return nil, err
  }

  // Set RaidName and RaidDifficulty
	tokens := strings.Split(maniestResponse.DisplayProperties.Name, ":")
	rawRaidName := strings.TrimSpace(tokens[0])
	rawRaidDifficulty := "NORMAL" // Default difficulty
	if len(tokens) > 1 {
		rawRaidDifficulty = strings.ToUpper(strings.TrimSpace(tokens[1]))
	}

	raidName, err := utils.ValidateRaidName(rawRaidName)
  if err != nil {
    return nil, err
  }
  *&entity.RaidName = raidName

  RaidDifficulty, err := utils.ValidateRaidDifficulty(rawRaidDifficulty)
  if err != nil {
    return nil, err
  }
  *&entity.RaidDifficulty = RaidDifficulty

  var playersGrouped = make(map[int64][]dto.PostGameCarnageReportEntry)
  for _, entry := range pgcr.Entries {
    val, ok := playersGrouped[entry.Player.DestinyUserInfo.MembershipId]
    if ok {
      playersGrouped[entry.Player.DestinyUserInfo.MembershipId] = append(val, entry) 
    } else {
      playersGrouped[entry.Player.DestinyUserInfo.MembershipId] = []dto.PostGameCarnageReportEntry{entry}
    }
  }

  // Process player information
  playerInformation, err := processPlayerInformation(playersGrouped, redisClient)
  if err != nil {
    return nil, err
  }

  *&entity.PlayerInformation = playerInformation
	return entity, nil
}

// Resolves if a raid was fresh or not, courtesy of @Newo
func resolveFromBeginning(pgcr *dto.PostGameCarnageReport, flawless bool) (*bool, error) {
	var result *bool = new(bool)

	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		return nil, err
	}

	if startTime.After(hauntedStart) || startTime.Equal(hauntedStart) {
		return &pgcr.ActivityWasStartedFromBeginning, nil
	} else if startTime.Before(beyondLightStart) {

		scourgeCriteria := pgcr.ActivityDetails.ActivityHash == sotpHash1 || pgcr.ActivityDetails.ActivityHash == sotpHash2
		leviathanCriteria := leviHashes[pgcr.ActivityDetails.ActivityHash]

		if scourgeCriteria {
			*result = pgcr.StartingPhaseIndex <= 1
			return result, nil
		} else if leviathanCriteria {
			*result = pgcr.StartingPhaseIndex == 0 || pgcr.StartingPhaseIndex == 2
			return result, nil
		} else {
			*result = pgcr.StartingPhaseIndex == 0
			return result, nil
		}
	} else if startTime.After(witchQueenStart) && (pgcr.ActivityWasStartedFromBeginning || flawless) {
		return &pgcr.ActivityWasStartedFromBeginning, nil
	}

	return result, nil
}

// Takes in a map of grouped up PGCR entries by players' membershipIds and returns an array of PlayerInformation structs
// Ensures that each player will have all their characters respectively
func processPlayerInformation(groups map[int64][]dto.PostGameCarnageReportEntry, client *redis.Client) ([]entity.PlayerInformation, error) {
  result := make([]entity.PlayerInformation, len(groups))
  for membershipId, value := range groups {
    playerInfo := entity.PlayerInformation{
      MembershipType: value[0].Player.DestinyUserInfo.MembershipType,
      MembershipId: value[0].Player.DestinyUserInfo.MembershipId,
      DisplayName: value[0].Player.DestinyUserInfo.DisplayName,
      GlobalDisplayName: value[0].Player.DestinyUserInfo.BungieGlobalDisplayName,
      GlobalDisplayNameCode: value[0].Player.DestinyUserInfo.BungieGlobalDisplayNameCode,
    }

    var playerCharacterInformation []entity.PlayerCharacterInformation
    for _, e := range value {
      characterInfo, err := createPlayerCharacter(&e, client) 
      if err != nil {
        log.Fatalf("There was an error create character information for player with Id [%d]\n", membershipId)
        return nil, err
      }
      playerCharacterInformation = append(playerCharacterInformation, *characterInfo) 
    }
    result = append(result, playerInfo)
  } 
  return result, nil
}

// Create an individual player character info struct based on a PGCR entry
// This utilizes Redis to fetch several pre-indexed manifest objects
// If querying Redis fails then this method return an error
func createPlayerCharacter(entry *dto.PostGameCarnageReportEntry, client *redis.Client) (*entity.PlayerCharacterInformation, error) {
  characterInfo := entity.PlayerCharacterInformation{
    ActivityCompleted: entry.Values["completed"].Basic.Value == 1.0,
    WeaponInformation: []entity.CharacterWeaponInformation{}, // empty just in case the player didn't do anything in the activity
  }
  rawClass, err := r.GetManifestEntity(client, entry.Player.CharacterClass) 
  if err != nil {
    return nil, err
  }
  class := utils.MapCharacterClass(rawClass.DisplayProperties.Name)

  rawGender, err := r.GetManifestEntity(client, strconv.FormatInt(entry.Player.GenderHash, 10))
  if err != nil {
    return nil, err
  }
  gender := utils.MapCharacterGender(rawGender.DisplayProperties.Name)

  rawRace, err := r.GetManifestEntity(client, strconv.FormatInt(entry.Player.RaceHash, 10))
  if err != nil {
    return nil, err
  }
  race := utils.MapCharacterRace(rawRace.DisplayProperties.Name)
  
  characterInfo.CharacterId = entry.CharacterId
  characterInfo.LightLevel = entry.Player.CharacterLevel
  characterInfo.CharacterClass = class
  characterInfo.CharacterRace = race
  characterInfo.CharacterGender = gender
  characterInfo.CharacterEmblem = entry.Player.EmblemHash
  characterInfo.TimePlayedSeconds = int(entry.Values["timePlayedSeconds"].Basic.Value)
  characterInfo.Kills = int(entry.Values["kills"].Basic.Value)
  characterInfo.Deaths = int(entry.Values["deaths"].Basic.Value)
  characterInfo.Assists = int(entry.Values["assists"].Basic.Value)
  characterInfo.Kda = entry.Values["killsDeathsAssits"].Basic.Value
  characterInfo.Kdr = entry.Values["killsDeathsRatio"].Basic.Value

  // Set weapon information
  if entry.Extended != nil {
    for _, weapon := range entry.Extended.Weapons {
      w := entity.CharacterWeaponInformation{
          WeaponHash: weapon.ReferenceId,
          Kills: int(weapon.Values["uniqueWeaponKills"].Basic.Value),
          PrecisionKills: int(weapon.Values["uniqueWeaponPrecisionKills"].Basic.Value),
          PrecisionRatio: weapon.Values["uniqueWeaponKillsPrecisionKills"].Basic.Value,
        }
      characterInfo.WeaponInformation = append(characterInfo.WeaponInformation, w)
    }

      // Set ability information
      abilityInfo := entity.CharacterAbilityInformation{
      GrenadeKills: int(entry.Extended.Abilities["weaponsKillsGrenade"].Basic.Value),
      MeleeKills: int(entry.Extended.Abilities["weaponsKillsMelee"].Basic.Value),
      SuperKills: int(entry.Extended.Abilities["weaponsKillsSuper"].Basic.Value),
    }
    characterInfo.AbilityInformation = abilityInfo
  }
    return &characterInfo, nil
}


