package service

import (
	"context"
	"database/sql"
	"fmt"
	"pgcr-processing-service/internal/model"
	"pgcr-processing-service/internal/repository"
	"strconv"
	"time"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/utils"
)

type PgcrService interface {
	ProcessPgcr(ppgcr types.ProcessedPostGameCarnageReport) error
}

type PgcrServiceImpl struct {
	Conn                          *sql.DB
	Redis                         ManifestClient
	RawPgcrRepository             repository.RawPgcrRepository
	PlayerRepository              repository.PlayerRepository
	RaidRepository                repository.RaidRepository
	InstanceActivityRepository    repository.InstanceActivityRepository
	PlayerRaidStatsRepository     repository.PlayerRaidStatsRepository
	WeaponRepository              repository.WeaponRepository
	InstanceWeaponStatsRepository repository.InstanceActivityWeaponStatsRepository
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

// Saves a processed pgcr to the Postgres DB
func (ps *PgcrServiceImpl) PersistProcessedPgcr(ppgcr types.ProcessedPostGameCarnageReport) error {
	tx, err := ps.Conn.Begin()
	if err != nil {
		return fmt.Errorf("Error opening transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	defer tx.Rollback()

	err = ps.addPlayerInfoToTx(tx, &ppgcr)
	if err != nil {
		return fmt.Errorf("Error adding players to transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	err = ps.addRaidInfoToTx(tx, &ppgcr)
	if err != nil {
		return fmt.Errorf("Error adding raids to transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	err = ps.addInstanceInfoToTx(tx, &ppgcr)
	if err != nil {
		return fmt.Errorf("Error adding instance activity to transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	err = ps.addWeaponInfoToTx(tx, &ppgcr)
	if err != nil {
		return fmt.Errorf("Error adding weapon stats to transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	err = ps.addOverallStatsToTx(tx, &ppgcr)
	if err != nil {
		return fmt.Errorf("Error adding instance stats to transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error commiting transaction for pgcr [%d]: %v", ppgcr.InstanceId, err)
	}

	return nil
}

// Adds info to the transaction regarding player specific stats
func (ps *PgcrServiceImpl) addOverallStatsToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	entities := mapOverallStats(ppgcr)

	for _, entity := range entities {
		_, err := ps.PlayerRaidStatsRepository.AddPlayerRaidStats(tx, entity)
		if err != nil {
			return fmt.Errorf("Error adding stats for player to transaction. MembershipId: [%d], Raid: [%s], Difficulty: [%s]: %v",
				entity.PlayerMembershipId, ppgcr.RaidName, ppgcr.RaidDifficulty, err)
		}
	}
	return nil
}

func mapOverallStats(ppgcr *types.ProcessedPostGameCarnageReport) []model.PlayerRaidStatsEntity {
	entities := []model.PlayerRaidStatsEntity{}
	for _, player := range ppgcr.PlayerInformation {
		for _, character := range player.PlayerCharacterInformation {
			statsEntity := model.PlayerRaidStatsEntity{
				RaidName:           ppgcr.RaidName,
				RaidDifficulty:     ppgcr.RaidDifficulty,
				PlayerMembershipId: player.MembershipId,
				Kills:              character.Kills,
				Deaths:             character.Deaths,
				Assists:            character.Assists,
				HoursPlayed:        character.TimePlayedSeconds,
				Flawless:           ppgcr.Flawless,
				Solo:               ppgcr.Solo,
				Duo:                ppgcr.Duo,
				Trio:               ppgcr.Trio,
			}

			clear := 0
			if character.ActivityCompleted {
				clear = 1
			}

			fullClear := 0
			if character.ActivityCompleted && ppgcr.FromBeginning {
				fullClear = 1
			}

			statsEntity.FullClears = fullClear
			statsEntity.Clears = clear
			entities = append(entities, statsEntity)
		}
	}
	return entities
}

func SetLowmanFlags(entity *model.PlayerRaidStatsEntity, flawless bool, playercount int) {
	entity.Flawless = flawless
	entity.Solo = playercount == 1 && entity.Clears == 1
	entity.Duo = playercount == 2 && entity.Clears == 1
	entity.Trio = playercount == 3 && entity.Clears == 1
	entity.SoloFlawless = flawless && playercount == 1 && entity.FullClears == 1
	entity.DuoFlawless = flawless && playercount == 2 && entity.FullClears == 1
	entity.TrioFlawless = flawless && playercount == 3 && entity.FullClears == 1
}

func (ps *PgcrServiceImpl) addInstanceInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, player := range ppgcr.PlayerInformation {
		for _, character := range player.PlayerCharacterInformation {
			entity := MapInstanceActivityEntity(ppgcr, &player, &character)
			_, err := ps.InstanceActivityRepository.AddInstanceActivity(tx, entity)
			if err != nil {
				return fmt.Errorf("Error adding instance activity [%s:%s] to transaction. Player MembershipId: [%d], CharacterId: [%d], Pgcr: [%d]",
					ppgcr.RaidName, ppgcr.RaidDifficulty, player.MembershipId, character.CharacterId, ppgcr.InstanceId)
			}
		}
	}
	return nil
}

func (ps *PgcrServiceImpl) addWeaponInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, player := range ppgcr.PlayerInformation {
		for _, character := range player.PlayerCharacterInformation {
			for _, weapon := range character.WeaponInformation {
				manifestEntity, err := ps.Redis.GetManifestEntity(context.Background(), strconv.FormatInt(weapon.WeaponHash, 10))
				if err != nil {
					return fmt.Errorf("Error while retrieving weapon with hash [%d] from Redis", weapon.WeaponHash)
				}
				weaponEntity := model.WeaponEntity{
					WeaponHash:          weapon.WeaponHash,
					WeaponIcon:          manifestEntity.DisplayProperties.Icon,
					WeaponName:          manifestEntity.DisplayProperties.Name,
					WeaponDamageType:    utils.GetDamageType(manifestEntity.EquippingBlock.AmmoType),
					WeaponEquipmentSlot: utils.GetEquippingSlot(int64(manifestEntity.EquippingBlock.EquipmentSlotTypeHash)),
				}

				_, err = ps.WeaponRepository.AddWeapon(tx, weaponEntity)
				if err != nil {
					return fmt.Errorf("Error while adding weapon with hash [%d] to the transaction: %v", weapon.WeaponHash, err)
				}

				raidActivityWeaponStatsEntity := model.InstanceWeaponStats{
					WeaponId:            weapon.WeaponHash,
					PlayerCharacterId:   character.CharacterId,
					InstanceId:          ppgcr.InstanceId,
					TotalKills:          weapon.Kills,
					TotalPrecisionKills: weapon.PrecisionKills,
					PrecisionRatio:      weapon.PrecisionRatio,
				}

				_, err = ps.InstanceWeaponStatsRepository.AddInstanceWeaponStats(tx, raidActivityWeaponStatsEntity)
				if err != nil {
					return fmt.Errorf("Error while adding weapon stats to transaction. WeaponId: [%d], CharacterId: [%d], Pgcr: [%d], MembershipId:MembershipType: [%d:%d]. %v",
						weapon.WeaponHash, character.CharacterId, ppgcr.InstanceId, player.MembershipId, player.MembershipType, err)
				}
			}
		}
	}
	return nil
}

func MapInstanceActivityEntity(ppgcr *types.ProcessedPostGameCarnageReport, player *types.PlayerData, character *types.PlayerCharacterInformation) model.InstanceActivityEntity {
	duration := ppgcr.EndTime.Sub(ppgcr.StartTime)
	return model.InstanceActivityEntity{
		InstanceId:         ppgcr.InstanceId,
		PlayerMembershipId: player.MembershipId,
		PlayerCharacterId:  character.CharacterId,
		CharacterEmblem:    character.CharacterEmblem,
		IsCompleted:        character.ActivityCompleted,
		Kills:              int32(character.Kills),
		Deaths:             int32(character.Deaths),
		Assists:            int32(character.Assists),
		MeleeKills:         character.AbilityInformation.MeleeKills,
		SuperKills:         character.AbilityInformation.SuperKills,
		GrenadeKills:       character.AbilityInformation.GrenadeKills,
		KillsDeathsAssists: character.Kda,
		KillsDeathsRatio:   character.Kdr,
		DurationSeconds:    int64(duration.Seconds()),
		TimeplayedSeconds:  int64(character.TimePlayedSeconds),
	}
}

func (ps *PgcrServiceImpl) addPlayerInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, playerInformation := range ppgcr.PlayerInformation {
		entity := MapPlayerEntity(playerInformation)
		_, err := ps.PlayerRepository.AddPlayer(tx, entity)
		if err != nil {
			return fmt.Errorf("Error adding player: %v", err)
		}
	}
	return nil
}

func MapPlayerEntity(playerData types.PlayerData) model.PlayerEntity {
	entity := model.PlayerEntity{
		MembershipId:   playerData.MembershipId,
		MembershipType: int32(playerData.MembershipType),
	}
	if playerData.GlobalDisplayName != "" {
		entity.DisplayName = playerData.GlobalDisplayName
	} else {
		entity.DisplayName = playerData.DisplayName
	}

	if playerData.GlobalDisplayNameCode != 0 {
		entity.DisplayNameCode = int32(playerData.GlobalDisplayNameCode)
	}

	characters := []model.PlayerCharacterEntity{}
	for _, playerCharacter := range playerData.PlayerCharacterInformation {
		entity := model.PlayerCharacterEntity{
			CharacterId:        playerCharacter.CharacterId,
			CharacterEmblem:    playerCharacter.CharacterEmblem,
			CharacterClass:     string(playerCharacter.CharacterClass),
			PlayerMembershipId: playerData.MembershipId,
		}
		characters = append(characters, entity)
	}
	entity.Characters = characters

	return entity
}

// Adds raids information to the SQL transaction
func (ps *PgcrServiceImpl) addRaidInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	entity := MapPgcrToRaidEntity(ppgcr)
	_, err := ps.RaidRepository.AddRaidInfo(tx, entity)
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
