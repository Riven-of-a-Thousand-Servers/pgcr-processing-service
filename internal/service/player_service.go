package service

import (
	"database/sql"
	"fmt"
	"pgcr-processing-service/internal/model"
	"pgcr-processing-service/internal/repository"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

type PlayerService struct {
	PlayerRepository repository.PlayerRepository
}

func (ps *PlayerService) ProcessPlayers(tx *sql.Tx, ppgcr *types.ProcessedPostGameCarnageReport) error {
	for _, playerInformation := range ppgcr.PlayerInformation {
		entity := PgcrToPlayerEntity(playerInformation)
		_, err := ps.PlayerRepository.AddPlayer(tx, entity)
		if err != nil {
			return fmt.Errorf("Error adding player: %v", err)
		}
	}
	return nil
}

func PgcrToPlayerEntity(playerData types.PlayerData) model.PlayerEntity {
	entity := model.PlayerEntity{
		MembershipId:   playerData.MembershipId,
		MembershipType: int32(playerData.MembershipType),
	}
	if playerData.GlobalDisplayName != "" {
		entity.DisplayName = playerData.GlobalDisplayName
		entity.DisplayNameCode = int32(playerData.GlobalDisplayNameCode)
	} else {
		entity.DisplayName = playerData.DisplayName
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
