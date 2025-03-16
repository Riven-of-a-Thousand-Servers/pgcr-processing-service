package mapper

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"rivenbot/internal/client"
	"rivenbot/internal/dto"
	"rivenbot/internal/model"
	"rivenbot/internal/utils"
)

type PgcrMapper struct {
	RedisClient client.RedisClient
}

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

// This method maps the PGCR into a pre-processed format thats more suitable for features
// Additionally, it compresses the raw PGCR fetched from Bungie and returns them if the compression is successful
func (p *PgcrMapper) Map(pgcr *dto.PostGameCarnageReport) ([]byte, *model.ProcessedPostGameCarnageReport, error) {
	processedPgcr, err := processPgcr(pgcr, p.RedisClient)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	compressed, err := utils.Compress(pgcr)
	if err != nil {
		return nil, nil, err
	}

	return compressed, processedPgcr, nil
}

func processPgcr(pgcr *dto.PostGameCarnageReport, redisClient client.RedisClient) (*model.ProcessedPostGameCarnageReport, error) {
	var entity model.ProcessedPostGameCarnageReport

	// Calculate start and end time
	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		log.Panicf("Something went wrong when parsing the period for PGCR [%s]: %v", pgcr.ActivityDetails.InstanceId, err)
		return nil, err
	}

	if len(pgcr.Entries) == 0 {
		error := fmt.Errorf("No entries in the PGCR [%s], unable to determine the end time of the activity", pgcr.ActivityDetails.InstanceId)
		return nil, error
	}

	activityDurationSeconds := int32(pgcr.Entries[0].Values["activityDurationSeconds"].Basic.Value)
	endTime := startTime.Add(time.Second * time.Duration(activityDurationSeconds))

	entity.StartTime = startTime
	entity.EndTime = endTime

	instanceId, err := strconv.ParseInt(pgcr.ActivityDetails.InstanceId, 10, 64)
	if err != nil {
		log.Panicf("Unable to convert instanceId [%s] to int64 for some reason?", pgcr.ActivityDetails.InstanceId)
		return nil, err
	}

	entity.InstanceId = instanceId
	entity.ActivityHash = pgcr.ActivityDetails.ActivityHash

	activitiyHash := pgcr.ActivityDetails.ActivityHash
	maniestResponse, err := redisClient.GetManifestEntity(context.Background(), strconv.Itoa(int(activitiyHash)))
	if err != nil {
		log.Panicf("Unable to find activity hash [%d] in Redis: %v", activitiyHash, err)
		return nil, err
	}

	raidName, raidDifficulty, err := utils.GetRaidAndDifficulty(maniestResponse.DisplayProperties.Name)
	if err != nil {
		log.Panic("Unable to parse activity raid name and raid difficulty")
		return nil, err
	}

	entity.RaidName = raidName
	entity.RaidDifficulty = raidDifficulty

	var playersGrouped = make(map[int64][]dto.PostGameCarnageReportEntry)
	for _, entry := range pgcr.Entries {
		membershipId, err := strconv.ParseInt(entry.Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			log.Panicf("Something went wrong when parsing membership ID [%s] to Int64", entry.Player.DestinyUserInfo.MembershipId)
		}
		val, ok := playersGrouped[membershipId]
		if ok {
			playersGrouped[membershipId] = append(val, entry)
		} else {
			playersGrouped[membershipId] = []dto.PostGameCarnageReportEntry{entry}
		}
	}

	// Process player information
	playerInformation, err := processPlayerInformation(playersGrouped)
	if err != nil {
		return nil, err
	}

	entity.PlayerInformation = playerInformation

	flawless := true

Outerloop:
	for _, players := range playersGrouped {
		for _, player := range players {
			if player.Values["deaths"].Basic.Value > 0 {
				flawless = false
				break Outerloop
			}
		}
	}

	fresh, err := resolveFromBeginning(pgcr, flawless)
	if err != nil {
		log.Panicf("Error parsing StartTime when determining if PGCR [%d] is fresh: %v", instanceId, err)
		return nil, err
	}

	trio := len(playersGrouped) == 3
	duo := len(playersGrouped) == 2
	solo := len(playersGrouped) == 1

	entity.Trio = trio
	entity.Duo = duo
	entity.Solo = solo
	entity.Flawless = flawless
	entity.FromBeginning = *fresh
	return &entity, nil
}

// Takes in a map of grouped up PGCR entries by players' membershipIds and returns an array of PlayerInformation structs
// Ensures that each player will have all their characters respectively
func processPlayerInformation(groups map[int64][]dto.PostGameCarnageReportEntry) ([]model.PlayerInformation, error) {
	result := []model.PlayerInformation{}
	for membershipId, entries := range groups {
		if len(entries) == 0 {
			log.Printf("Player with membershipId [%d] has no entries, skipping\n", membershipId)
			continue
		}

		playerInfo := model.PlayerInformation{
			MembershipId:          membershipId,
			MembershipType:        entries[0].Player.DestinyUserInfo.MembershipType,
			DisplayName:           entries[0].Player.DestinyUserInfo.DisplayName,
			GlobalDisplayName:     entries[0].Player.DestinyUserInfo.BungieGlobalDisplayName,
			GlobalDisplayNameCode: entries[0].Player.DestinyUserInfo.BungieGlobalDisplayNameCode,
		}

		var characters []model.PlayerCharacterInformation
		for _, e := range entries {
			characterInfo, err := createPlayerCharacter(&e)
			if err != nil {
				log.Panicf("There was an error create character information for player with Id [%d]: %v\n", membershipId, err)
				return nil, err
			}
			characters = append(characters, *characterInfo)
		}

		playerInfo.PlayerCharacterInformation = characters
		result = append(result, playerInfo)
	}

	return result, nil
}

// Create an individual player character info struct based on a PGCR entry
// This utilizes Redis to fetch several pre-indexed manifest objects
// If querying Redis fails then this method return an error
func createPlayerCharacter(entry *dto.PostGameCarnageReportEntry) (*model.PlayerCharacterInformation, error) {
	characterInfo := model.PlayerCharacterInformation{
		ActivityCompleted: entry.Values["completed"].Basic.Value == 1.0,
		WeaponInformation: []model.CharacterWeaponInformation{}, // empty just in case the player didn't do anything in the activity
	}

	class := model.CharacterClass(entry.Player.CharacterClass)

	characterId, err := strconv.ParseInt(entry.CharacterId, 10, 64)
	if err != nil {
		log.Panicf("Unable to parse character Id [%s] to int64", entry.CharacterId)
		return nil, err
	}

	characterInfo.CharacterId = characterId
	characterInfo.LightLevel = entry.Player.LightLevel
	characterInfo.CharacterClass = class
	characterInfo.CharacterEmblem = entry.Player.EmblemHash
	characterInfo.TimePlayedSeconds = int(entry.Values["timePlayedSeconds"].Basic.Value)
	characterInfo.Kills = int(entry.Values["kills"].Basic.Value)
	characterInfo.Deaths = int(entry.Values["deaths"].Basic.Value)
	characterInfo.Assists = int(entry.Values["assists"].Basic.Value)
	characterInfo.Kda = entry.Values["killsDeathsAssists"].Basic.Value
	characterInfo.Kdr = entry.Values["killsDeathsRatio"].Basic.Value

	// Set weapon information
	if entry.Extended != nil {
		for _, weapon := range entry.Extended.Weapons {
			w := model.CharacterWeaponInformation{
				WeaponHash:     weapon.ReferenceId,
				Kills:          int(weapon.Values["uniqueWeaponKills"].Basic.Value),
				PrecisionKills: int(weapon.Values["uniqueWeaponPrecisionKills"].Basic.Value),
				PrecisionRatio: weapon.Values["uniqueWeaponKillsPrecisionKills"].Basic.Value,
			}
			characterInfo.WeaponInformation = append(characterInfo.WeaponInformation, w)
		}

		// Set ability information
		abilityInfo := model.CharacterAbilityInformation{
			GrenadeKills: int(entry.Extended.Abilities["weaponKillsGrenade"].Basic.Value),
			MeleeKills:   int(entry.Extended.Abilities["weaponKillsMelee"].Basic.Value),
			SuperKills:   int(entry.Extended.Abilities["weaponKillsSuper"].Basic.Value),
		}
		characterInfo.AbilityInformation = abilityInfo
	}
	return &characterInfo, nil
}

func GroupCharacters(entries []dto.PostGameCarnageReportEntry) (map[int64][]dto.PostGameCarnageReportEntry, error) {
	var playersGrouped = make(map[int64][]dto.PostGameCarnageReportEntry)
	for _, entry := range entries {
		membershipId, err := strconv.ParseInt(entry.Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			log.Panicf("Something went wrong when parsing membership ID [%s] to Int64", entry.Player.DestinyUserInfo.MembershipId)
			return nil, err
		}
		val, ok := playersGrouped[membershipId]
		if ok {
			playersGrouped[membershipId] = append(val, entry)
		} else {
			playersGrouped[membershipId] = []dto.PostGameCarnageReportEntry{entry}
		}
	}
	return playersGrouped, nil
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

		isScourge := pgcr.ActivityDetails.ActivityHash == sotpHash1 || pgcr.ActivityDetails.ActivityHash == sotpHash2
		isLeviathan := leviHashes[pgcr.ActivityDetails.ActivityHash]

		if isScourge {
			*result = pgcr.StartingPhaseIndex <= 1
			return result, nil
		} else if isLeviathan {
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
