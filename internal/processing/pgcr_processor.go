package processing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Riven-of-a-thousand-servers/commons"
	"pgcr-processing-service/internal/db"
	"pgcr-processing-service/internal/mapper"
	"pgcr-processing-service/internal/rabbitmq"
	"pgcr-processing-service/internal/redis"
)

type Processor interface {
	DoPgcr(types.ProcessedPostGameCarnageReport) error
}

type PgcrProcessor struct {
	Querier    db.Querier
	RabbitMq   *rabbitmq.RabbitMQ
	PgcrMapper mapper.PgcrMapper
	Redis      redis.Service
}

func NewPgcrProcessor(querier db.Querier, rabbitmq *rabbitmq.RabbitMQ, redis redis.Service) *PgcrProcessor {
	return &PgcrProcessor{
		Querier:  querier,
		RabbitMq: rabbitmq,
		Redis:    redis,
	}
}

// Saves a processed pgcr to the Postgres DB
func (p *PgcrProcessor) DoPgcr(ctx context.Context, pgcr types.ProcessedPostGameCarnageReport, b []byte) error {
	instance := db.CreateInstanceParams{
		ID:              pgcr.InstanceId,
		ActivityHash:    pgcr.ActivityHash,
		IsFresh:         pgcr.FromBeginning,
		PlayerCount:     int32(len(pgcr.PlayerInformation)),
		StartTime:       pgcr.StartTime,
		EndTime:         pgcr.EndTime,
		DurationSeconds: int32(pgcr.EndTime.Sub(pgcr.StartTime).Seconds()),
	}

	rawPgcr := db.CreatePgcrParams{
		InstanceID: pgcr.InstanceId,
		Blob:       b,
	}

	// Player
	var destinyPlayer []db.CreateDestinyPlayerParams
	for _, pi := range pgcr.PlayerInformation {
		player := db.CreateDestinyPlayerParams{
			MembershipID:   pi.MembershipId,
			MembershipType: int32(pi.MembershipType),
		}

		if pi.GlobalDisplayName != "" {
			player.DisplayName = sql.NullString{String: pi.GlobalDisplayName}
		} else {
			player.DisplayName = sql.NullString{String: pi.DisplayName}
		}

		if pi.GlobalDisplayNameCode != 0 {
			player.GlobalDisplayNameCode = sql.NullInt32{
				Int32: int32(pi.GlobalDisplayNameCode),
			}
		}

		// InstancePlayer
		instancePlayer := db.CreateInstancePlayerParams{
			InstanceID:   pgcr.InstanceId,
			MembershipID: pi.MembershipId,
		}

		// InstanceCharacter
		var instancePlayerCharacters []db.CreateInstanceCharacterParams
		for _, ci := range pi.PlayerCharacterInformation {
			instanceCharacter := db.CreateInstanceCharacterParams{
				InstanceID:   pgcr.InstanceId,
				MembershipID: pi.MembershipId,
				CharacterID:  ci.CharacterId,
				EmblemHash:   ci.CharacterEmblem,
				Completed:    ci.ActivityCompleted,
				Kills:        int32(ci.Kills),
				Deaths:       int32(ci.Deaths),
				Assists:      int32(ci.Assists),
				Kda:          strconv.FormatFloat(float64(ci.Kda), 'f', -1, 64),
				Kdr:          strconv.FormatFloat(float64(ci.Kdr), 'f', -1, 64),
				Efficiency:   ci.Efficiency,
				SuperKills:   int32(ci.AbilityInformation.SuperKills),
				GrenadeKills: int32(ci.AbilityInformation.GrenadeKills),
				MeleeKills:   int32(ci.AbilityInformation.MeleeKills),
			}
		}

		destinyPlayer = append(destinyPlayer, player)
	}

	// InstancePlayer
	for _, pi := range pgcr.PlayerInformation {
		if err := p.Querier.CreateInstancePlayer(ctx, db.CreateInstancePlayerParams{
			InstanceID:   pgcr.InstanceId,
			MembershipID: pi.MembershipId,
		}); err != nil {
		}
	}

	return nil
	// err = p.addPlayerInfoToTx(tx, &pgcr)
	// if err != nil {
	// 	return fmt.Errorf("Error adding players to transaction for pgcr [%d]: %v", pgcr.InstanceId, err)
	// }
	//
	// err = p.addRaidInfoToTx(tx, &pgcr)
	// if err != nil {
	// 	return fmt.Errorf("Error adding raids to transaction for pgcr [%d]: %v", pgcr.InstanceId, err)
	// }
	//
	// err = p.addInstanceInfoToTx(tx, &pgcr)
	// if err != nil {
	// 	return fmt.Errorf("Error adding instance activity to transaction for pgcr [%d]: %v", pgcr.InstanceId, err)
	// }
	//
	// err = p.addWeaponInfoToTx(tx, &pgcr)
	// if err != nil {
	// 	return fmt.Errorf("Error adding weapon stats to transaction for pgcr [%d]: %v", pgcr.InstanceId, err)
	// }
	//
	// err = p.addOverallStatsToTx(tx, &pgcr)
	// if err != nil {
	// 	return fmt.Errorf("Error adding instance stats to transaction for pgcr [%d]: %v", pgcr.InstanceId, err)
	// }
	//
	// err = tx.Commit()
	// if err != nil {
	// 	return fmt.Errorf("Error commiting transaction for pgcr [%d]: %v", pgcr.InstanceId, err)
	// }
	//
	// return nil
}

func (p *PgcrProcessor) Consume(ctx context.Context) error {
	deliveries, err := p.RabbitMq.Consumer("pgcr_consumer")
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		var pgcr types.PostGameCarnageReportResponse
		err := json.Unmarshal(delivery.Body, &pgcr)
		if err != nil {
			slog.Error("Error unmarshalling body from message", "Error", err)
			return fmt.Errorf("Error unmarshalling body from message: %v", err)
		}

		instanceId := pgcr.Response.ActivityDetails.InstanceId
		slog.Info("Processing pgcr", "InstanceId", instanceId)
		b, ppgcr, err := p.PgcrMapper.ToProcessedPgcr(&pgcr.Response)
		if err != nil {
			return fmt.Errorf("Error mapping pgcr [%s] to a processed pgcr: %v", instanceId, err)
		}
		err = p.DoPgcr(ctx, *ppgcr, b)
		if err != nil {
			return fmt.Errorf("Error processing ppgcr [%d] into database tables: %v", ppgcr.InstanceId, err)
		}
		slog.Info("Finished processing pgcr", "InstanceId", instanceId)
	}
	return nil
}

// Adds info to the transaction regarding player specific stats
func (p *Processor) addOverallStatsToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	entities := mapOverallStats(ppgcr)

	for _, entity := range entities {
		_, err := p.PlayerRaidStatsRepository.AddPlayerRaidStats(tx, entity)
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
				SoloFlawless:       ppgcr.Solo && ppgcr.Flawless,
				TrioFlawless:       ppgcr.Trio && ppgcr.Flawless,
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

func (p *Processor) addInstanceInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, player := range ppgcr.PlayerInformation {
		for _, character := range player.PlayerCharacterInformation {
			entity := MapInstanceActivityEntity(ppgcr, &player, &character)
			_, err := p.InstanceActivityRepository.AddInstanceActivity(tx, entity)
			if err != nil {
				return fmt.Errorf("Error adding instance activity [%s:%s] to transaction. Player MembershipId: [%d], CharacterId: [%d], Pgcr: [%d]",
					ppgcr.RaidName, ppgcr.RaidDifficulty, player.MembershipId, character.CharacterId, ppgcr.InstanceId)
			}
		}
	}
	return nil
}

func (p *Processor) addWeaponInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, player := range ppgcr.PlayerInformation {
		for _, character := range player.PlayerCharacterInformation {
			for _, weapon := range character.WeaponInformation {
				manifestEntity, err := p.Redis.GetManifestEntity(context.Background(), strconv.FormatInt(weapon.WeaponHash, 10))
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

				_, err = p.WeaponRepository.AddWeapon(tx, weaponEntity)
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

				_, err = p.InstanceWeaponStatsRepository.AddInstanceWeaponStats(tx, raidActivityWeaponStatsEntity)
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

func (p *Processor) addPlayerInfoToTx(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, playerInformation := range ppgcr.PlayerInformation {
		entity := MapPlayerEntity(playerInformation)
		_, err := p.PlayerRepository.AddPlayer(tx, entity)
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
