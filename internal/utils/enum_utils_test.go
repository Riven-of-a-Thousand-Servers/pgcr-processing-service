package utils

import (
	"fmt"
	"rivenbot/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

var enumUtilTests = map[string]struct {
	input      string
	raid       model.RaidName
	difficulty model.RaidDifficulty
}{
	"leviathan_normal": {
		input:      "Leviathan: Normal",
		raid:       model.LEVIATHAN,
		difficulty: model.NORMAL,
	},
	"leviathan_prestige": {
		input:      "Leviathan: Prestige",
		raid:       model.LEVIATHAN,
		difficulty: model.PRESTIGE,
	},
	"leviathan_guided_games": {
		input:      "Leviathan: Guided Games",
		raid:       model.LEVIATHAN,
		difficulty: model.GUIDED_GAMES,
	},
	"eater_of_worlds_normal": {
		input:      "Leviathan, Eater of Worlds: Normal",
		raid:       model.EATER_OF_WORLDS,
		difficulty: model.NORMAL,
	},
	"eater_of_worlds_prestige": {
		input:      "Leviathan, Eater of Worlds: Prestige",
		raid:       model.EATER_OF_WORLDS,
		difficulty: model.PRESTIGE,
	},
	"eater_of_worlds_guided_games": {
		input:      "Leviathan, Eater of Worlds: Guided Games",
		raid:       model.EATER_OF_WORLDS,
		difficulty: model.GUIDED_GAMES,
	},
	"spire_of_stars_normal": {
		input:      "Leviathan, Spire of Stars: Normal",
		raid:       model.SPIRE_OF_STARS,
		difficulty: model.NORMAL,
	},
	"spire_of_stars_prestige": {
		input:      "Leviathan, Spire of Stars: Prestige",
		raid:       model.SPIRE_OF_STARS,
		difficulty: model.PRESTIGE,
	},
	"spire_of_stars_guided_games": {
		input:      "Leviathan, Spire of Stars: Guided Games",
		raid:       model.SPIRE_OF_STARS,
		difficulty: model.GUIDED_GAMES,
	},
	"scourge_of_the_past": {
		input:      "Scourge of the Past",
		raid:       model.SCOURGE_OF_THE_PAST,
		difficulty: model.NORMAL,
	},
	"last_wish": {
		input:      "Last Wish",
		raid:       model.LAST_WISH,
		difficulty: model.NORMAL,
	},
	"crown_of_sorrow": {
		input:      "Crown of Sorrow",
		raid:       model.CROWN_OF_SORROW,
		difficulty: model.NORMAL,
	},
	"garden_of_salvation": {
		input:      "Garden of Salvation",
		raid:       model.GARDEN_OF_SALVATION,
		difficulty: model.NORMAL,
	},
	"deep_stone_crypt": {
		input:      "Deep Stone Crypt",
		raid:       model.DEEP_STONE_CRYPT,
		difficulty: model.NORMAL,
	},
	"vault_of_glass_normal": {
		input:      "Vault of Glass: Standard",
		raid:       model.VAULT_OF_GLASS,
		difficulty: model.NORMAL,
	},
	"vault_of_glass_master": {
		input:      "Vault of Glass: Master",
		raid:       model.VAULT_OF_GLASS,
		difficulty: model.MASTER,
	},
	"vault_of_glass_challenge": {
		input:      "Vault of Glass: Challenge Mode",
		raid:       model.VAULT_OF_GLASS,
		difficulty: model.CHALLENGE_MODE,
	},
	"vow_of_the_disciple_normal": {
		input:      "Vow of the Disciple: Standard",
		raid:       model.VOW_OF_THE_DISCIPLE,
		difficulty: model.NORMAL,
	},
	"vow_of_the_disciple_master": {
		input:      "Vow of the Disciple: Master",
		raid:       model.VOW_OF_THE_DISCIPLE,
		difficulty: model.MASTER,
	},
	"kings_fall_challenge_mode": {
		input:      "King's Fall: Expert",
		raid:       model.KINGS_FALL,
		difficulty: model.CHALLENGE_MODE,
	},
	"kings_fall_normal": {
		input:      "King's Fall: Standard",
		raid:       model.KINGS_FALL,
		difficulty: model.NORMAL,
	},
	"kings_fall_master": {
		input:      "King's Fall: Master",
		raid:       model.KINGS_FALL,
		difficulty: model.MASTER,
	},
	"root_of_nightmares_normal": {
		input:      "Root of Nightmares: Standard",
		raid:       model.ROOT_OF_NIGHTMARES,
		difficulty: model.NORMAL,
	},
	"root_of_nightmares_master": {
		input:      "Root of Nightmares: Master",
		raid:       model.ROOT_OF_NIGHTMARES,
		difficulty: model.MASTER,
	},
	"crotas_end_challenge_mode": {
		input:      "Crota's End: Legend",
		raid:       model.CROTAS_END,
		difficulty: model.CHALLENGE_MODE,
	},
	"crotas_end_normal": {
		input:      "Crota's End: Normal",
		raid:       model.CROTAS_END,
		difficulty: model.NORMAL,
	},
	"crotas_end_master": {
		input:      "Crota's End: Master",
		raid:       model.CROTAS_END,
		difficulty: model.MASTER,
	},
	"salvations_edge_normal": {
		input:      "Salvation's Edge: Standard",
		raid:       model.SALVATIONS_EDGE,
		difficulty: model.NORMAL,
	},
	// This one is a contest mode clear
	"salvations_edge_normal_2.0": {
		input:      "Salvation's Edge",
		raid:       model.SALVATIONS_EDGE,
		difficulty: model.NORMAL,
	},
	"salvations_edge_master": {
		input:      "Salvation's Edge: Master",
		raid:       model.SALVATIONS_EDGE,
		difficulty: model.MASTER,
	},
}

func TestGetRaidAndDifficulty_Success(t *testing.T) {
	// given: A string representing a raid
	for test, params := range enumUtilTests {
		t.Run(test, func(t *testing.T) {
			// when: GetRaidAndDifficulty is called
			raid, difficulty, err := GetRaidAndDifficulty(params.input)

			if err != nil {
				t.Fatalf("Expected no errors, found one: %v", err)
			}

			assert := assert.New(t)
			// then: Raid and Difficulty values are correct
			assert.Equal(params.raid, raid, fmt.Sprintf("Raid should be [%s]", params.raid))
			assert.Equal(params.difficulty, difficulty, fmt.Sprintf("Difficulty should be [%s]", params.difficulty))
		})
	}
}
