package pgcr

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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
// Subsequently it compresses it using Gzip
func ProcessPgcr(pgcr *dto.PostGameCarnageReport, db *sql.DB) {
	processedPgcr, err := createProcessedPgcr(pgcr)
	if err != nil {
		log.Fatal(err)
	}
}

func createProcessedPgcr(pgcr *dto.PostGameCarnageReport) (*entity.ProcessedPostGameCarnageReport, error) {
	entity := new(entity.ProcessedPostGameCarnageReport)
	startTime, err := time.Parse(time.RFC3339, pgcr.Period)
	if err != nil {
		return nil, err
	}

	if len(pgcr.Entries) == 0 {
		error := fmt.Errorf("No entries in the PostGameCarnageReport, unable to determine the end time of the activity")
		return nil, error
	}

	activityDurationSeconds := int32(pgcr.Entries[0].Values["activityDurationSeconds"].Value)
	endTime := startTime.Add(time.Second * time.Duration(activityDurationSeconds))

	*&entity.StartTime = startTime
	*&entity.EndTime = endTime
	*&entity.InstanceId = pgcr.ActivityDetails.InstanceId

	tokens := strings.Split(manifestResponseName, ":")
	rawRaidName := strings.TrimSpace(tokens[0])
	rawRaidDifficulty := "NORMAL" // Default difficulty
	if len(tokens) > 1 {
		rawRaidDifficulty = strings.ToUpper(strings.TrimSpace(tokens[1]))
	}

	*&entity.RaidName, err := validateRaidName(rawRaidName)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

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

func validateRaidName(name string) (entity.RaidName, error) {
	switch entity.RaidName(strings.ToUpper(strings.TrimSpace(name))) {
	case entity.SalvationsEdge:
		return entity.SalvationsEdge, nil
	case entity.CrotasEnd:
		return entity.CrotasEnd, nil
	case entity.RootOfNightmares:
		return entity.RootOfNightmares, nil
	case entity.KingsFall:
		return entity.KingsFall, nil
	case entity.VowOfTheDisciple:
		return entity.VowOfTheDisciple, nil
	case entity.VaultOfGlass:
		return entity.VaultOfGlass, nil
	case entity.DeepStoneCrypt:
		return entity.DeepStoneCrypt, nil
	case entity.GardenOfSalvation:
		return entity.GardenOfSalvation, nil
	case entity.LeviathanCrownOfSorrow:
		return entity.LeviathanCrownOfSorrow, nil
	case entity.LastWish:
		return entity.LastWish, nil
	case entity.LeviathanSpireOfStars:
		return entity.LeviathanSpireOfStars, nil
	case entity.LeviathanEaterOfWorlds:
		return entity.LeviathanEaterOfWorlds, nil
	case entity.Leviathan:
		return entity.Leviathan, nil
	case entity.ScourgeOfThePast:
		return entity.ScourgeOfThePast, nil
	default:
		return "", errors.New("invalid raid name")
	}
}

func validateRaidDifficulty(name string) (entity.RaidDifficulty, error) {
	switch entity.RaidDifficulty(strings.ToUpper(strings.TrimSpace(name))) {
	case entity.Normal:
		return entity.Normal, nil
	case entity.Prestige:
		return entity.Prestige, nil
	case entity.Master:
		return entity.Master, nil
	default:
		return "", errors.New("invalid raid difficulty")
	}
}
