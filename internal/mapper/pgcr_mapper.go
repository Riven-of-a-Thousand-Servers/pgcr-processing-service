package mapper

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"pgcr-processing-service/internal/compress"
	"pgcr-processing-service/internal/redis"
	"pgcr-processing-service/internal/types/pgcr"
	"pgcr-processing-service/internal/types/rabbitmq"
	"pgcr-processing-service/internal/utils"
)

type PgcrMapper struct {
	ManifestClient redis.Service
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
func (p *PgcrMapper) ToProcessedPgcr(pgcr *pgcr.PostGameCarnageReport) ([]byte, *rabbitmq.ProcessedPostGameCarnageReport, error) {
	processedPgcr, err := processPgcr(pgcr, p.ManifestClient)
	if err != nil {
		return nil, nil, err
	}

	compressed, err := compress.Gzip(pgcr)
	if err != nil {
		return nil, nil, err
	}

	return compressed, processedPgcr, nil
}

func processPgcr(report *pgcr.PostGameCarnageReport, redisService redis.Service) (*rabbitmq.ProcessedPostGameCarnageReport, error) {
	var entity rabbitmq.ProcessedPostGameCarnageReport

	// Calculate start and end time
	startTime, err := time.Parse(time.RFC3339, report.Period)
	if err != nil {
		slog.Error("Something went wrong when parsing the period for PGCR", "InstanceId", report.ActivityDetails.InstanceId, "Error", err)
		return nil, err
	}

	if len(report.Entries) == 0 {
		error := fmt.Errorf("No entries in the PGCR [%s], unable to determine the end time of the activity", report.ActivityDetails.InstanceId)
		return nil, error
	}

	activityDurationSeconds := int32(report.Entries[0].Values["activityDurationSeconds"].Basic.Value)
	endTime := startTime.Add(time.Second * time.Duration(activityDurationSeconds))

	entity.StartTime = startTime
	entity.EndTime = endTime

	instanceId, err := strconv.ParseInt(report.ActivityDetails.InstanceId, 10, 64)
	if err != nil {
		slog.Error("Unable to convert instanceIdto int64 for some reason?", "InstanceId", report.ActivityDetails.InstanceId)
		return nil, err
	}

	entity.InstanceId = instanceId
	entity.ActivityHash = report.ActivityDetails.ActivityHash

	activitiyHash := report.ActivityDetails.ActivityHash
	manifestResponse, err := redisService.GetManifestEntity(context.Background(), strconv.Itoa(int(activitiyHash)))
	if err != nil {
		slog.Error("Unable to find activity hash in Redis", "ActivityHash", activitiyHash, "Error", err)
		return nil, err
	}

	raidName, raidDifficulty, err := utils.GetRaidAndDifficulty(manifestResponse.DisplayProperties.Name)
	if err != nil {
		slog.Error("Unable to parse activity raid name and raid difficulty")
		return nil, err
	}

	entity.RaidName = raidName
	entity.RaidDifficulty = raidDifficulty

	playersGrouped := make(map[int64][]pgcr.PostGameCarnageReportEntry)
	for _, entry := range report.Entries {
		membershipId, err := strconv.ParseInt(entry.Player.DestinyUserInfo.MembershipId, 10, 64)
		if err != nil {
			slog.Error("Something went wrong when parsing membership ID to Int64", "MembershipId", entry.Player.DestinyUserInfo.MembershipId)
			return nil, err
		}
		val, ok := playersGrouped[membershipId]
		if ok {
			playersGrouped[membershipId] = append(val, entry)
		} else {
			playersGrouped[membershipId] = []pgcr.PostGameCarnageReportEntry{entry}
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

	fresh, err := resolveFromBeginning(report, flawless)
	if err != nil {
		slog.Error("Error parsing StartTime when determining if PGCRis fresh", "InstanceId", instanceId, "Error", err)
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
func processPlayerInformation(groups map[int64][]pgcr.PostGameCarnageReportEntry) ([]rabbitmq.PlayerData, error) {
	result := []rabbitmq.PlayerData{}
	for membershipId, entries := range groups {
		if len(entries) == 0 {
			slog.Info("Player with membershipId has no entries, skipping", "MembershipId", membershipId)
			continue
		}

		playerInfo := rabbitmq.PlayerData{
			MembershipId:          membershipId,
			MembershipType:        entries[0].Player.DestinyUserInfo.MembershipType,
			DisplayName:           entries[0].Player.DestinyUserInfo.DisplayName,
			IsPublic:              entries[0].Player.DestinyUserInfo.IsPublic,
			IconPath:              entries[0].Player.DestinyUserInfo.IconPath,
			GlobalDisplayName:     entries[0].Player.DestinyUserInfo.BungieGlobalDisplayName,
			GlobalDisplayNameCode: entries[0].Player.DestinyUserInfo.BungieGlobalDisplayNameCode,
		}

		for _, e := range entries {
			characterInfo, err := createPlayerCharacter(&e)
			if err != nil {
				slog.Error("There was an error create character information for player with Id", "MembershipId", membershipId, "Error", err)
				return nil, err
			}
			playerInfo.PlayerCharacterInformation = append(playerInfo.PlayerCharacterInformation, *characterInfo)
		}

		for _, c := range playerInfo.PlayerCharacterInformation {
			if c.ActivityCompleted {
				playerInfo.Completed = true
				break
			}
		}

		totalTimePlayed := 0
		for _, c := range playerInfo.PlayerCharacterInformation {
			totalTimePlayed += c.TimePlayedSeconds
		}
		playerInfo.TimePlayedSeconds = int32(totalTimePlayed)
		result = append(result, playerInfo)
	}

	return result, nil
}

// Create an individual player character info struct based on a PGCR entry
// This utilizes Redis to fetch several pre-indexed manifest objects
// If querying Redis fails then this method return an error
func createPlayerCharacter(entry *pgcr.PostGameCarnageReportEntry) (*rabbitmq.PlayerCharacterInformation, error) {
	characterInfo := rabbitmq.PlayerCharacterInformation{
		ActivityCompleted: entry.Values["completed"].Basic.Value == 1.0,
		WeaponInformation: []rabbitmq.CharacterWeaponInformation{}, // empty just in case the player didn't do anything in the activity
	}

	class := rabbitmq.CharacterClass(entry.Player.CharacterClass)

	characterId, err := strconv.ParseInt(entry.CharacterId, 10, 64)
	if err != nil {
		slog.Error("Unable to parse character Id to int64", "CharacterId", entry.CharacterId)
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
			w := rabbitmq.CharacterWeaponInformation{
				WeaponHash:     weapon.ReferenceId,
				Kills:          int(weapon.Values["uniqueWeaponKills"].Basic.Value),
				PrecisionKills: int(weapon.Values["uniqueWeaponPrecisionKills"].Basic.Value),
				PrecisionRatio: weapon.Values["uniqueWeaponKillsPrecisionKills"].Basic.Value,
			}
			characterInfo.WeaponInformation = append(characterInfo.WeaponInformation, w)
		}

		// Set ability information
		abilityInfo := rabbitmq.CharacterAbilityInformation{
			GrenadeKills: int(entry.Extended.Abilities["weaponKillsGrenade"].Basic.Value),
			MeleeKills:   int(entry.Extended.Abilities["weaponKillsMelee"].Basic.Value),
			SuperKills:   int(entry.Extended.Abilities["weaponKillsSuper"].Basic.Value),
		}
		characterInfo.AbilityInformation = abilityInfo
	}
	return &characterInfo, nil
}

// Resolves if a raid was fresh or not, courtesy of @Newo
func resolveFromBeginning(pgcr *pgcr.PostGameCarnageReport, flawless bool) (*bool, error) {
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
