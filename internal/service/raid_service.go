package service

import (
	"database/sql"
	"fmt"
	"pgcr-processing-service/internal/model"
	"pgcr-processing-service/internal/repository"
	"time"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

type RaidService struct {
	RaidRepository repository.RaidRepository
}

var (
	pst = time.FixedZone("PST", -8*60*60)
	pdt = time.FixedZone("PDT", -7*60*60)
)

var activeRaids = map[types.RaidName]bool{
	types.SALVATIONS_EDGE:     true,
	types.CROTAS_END:          true,
	types.ROOT_OF_NIGHTMARES:  true,
	types.KINGS_FALL:          true,
	types.VOW_OF_THE_DISCIPLE: true,
	types.VAULT_OF_GLASS:      true,
	types.DEEP_STONE_CRYPT:    true,
	types.GARDEN_OF_SALVATION: true,
	types.LAST_WISH:           true,
	types.CROWN_OF_SORROW:     false,
	types.SCOURGE_OF_THE_PAST: false,
	types.SPIRE_OF_STARS:      false,
	types.EATER_OF_WORLDS:     false,
	types.LEVIATHAN:           false,
}

var raidReleaseDates = map[types.RaidName]time.Time{
	types.SALVATIONS_EDGE:     time.Date(2024, time.June, 7, 9, 0, 0, 0, pdt),
	types.CROTAS_END:          time.Date(2023, time.September, 1, 9, 0, 0, 0, pdt),
	types.ROOT_OF_NIGHTMARES:  time.Date(2023, time.March, 10, 9, 0, 0, 0, pst),
	types.KINGS_FALL:          time.Date(2022, time.August, 26, 9, 0, 0, 0, pdt),
	types.VOW_OF_THE_DISCIPLE: time.Date(2022, time.March, 5, 9, 0, 0, 0, pst),
	types.VAULT_OF_GLASS:      time.Date(2021, time.May, 22, 9, 0, 0, 0, pdt),
	types.DEEP_STONE_CRYPT:    time.Date(2020, time.November, 21, 9, 0, 0, 0, pst),
	types.GARDEN_OF_SALVATION: time.Date(2019, time.October, 5, 9, 0, 0, 0, pdt),
	types.CROWN_OF_SORROW:     time.Date(2019, time.June, 4, 9, 0, 0, 0, pdt),
	types.LAST_WISH:           time.Date(2018, time.September, 14, 9, 0, 0, 0, pdt),
	types.SCOURGE_OF_THE_PAST: time.Date(2018, time.December, 7, 9, 0, 0, 0, pst),
	types.SPIRE_OF_STARS:      time.Date(2018, time.May, 8, 9, 0, 0, 0, pdt),
	types.EATER_OF_WORLDS:     time.Date(2017, time.December, 6, 9, 0, 0, 0, pst),
	types.LEVIATHAN:           time.Date(2017, time.September, 13, 9, 0, 0, 0, pdt),
}

func (rs *RaidService) ProcessRaid(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	entity := MapPgcrToRaidEntity(ppgcr)
	_, err := rs.RaidRepository.AddRaidInfo(tx, entity)
	if err != nil {
		return fmt.Errorf("Error while adding raid to DB: %v", err)
	}
	return nil
}

func MapPgcrToRaidEntity(ppgcr *types.ProcessedPostGameCarnageReport) model.RaidEntity {
	return model.RaidEntity{
		RaidName:       string(ppgcr.RaidName),
		RaidDifficulty: string(ppgcr.RaidDifficulty),
		RaidHash:       ppgcr.ActivityHash,
		IsActive:       activeRaids[ppgcr.RaidName],
		ReleaseDate:    raidReleaseDates[ppgcr.RaidName],
	}
}
