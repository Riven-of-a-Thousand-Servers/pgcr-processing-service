package utils

import (
	"errors"
	"fmt"
	"log"
	"rivenbot/internal/model"
	"strings"
	"sync"
)

var reverseClassLabels map[string]model.CharacterClass
var characterClassLabels = map[model.CharacterClass]string{
	model.TITAN:   "Titan",
	model.WARLOCK: "Warlock",
	model.HUNTER:  "Hunter",
}

func ClassLabel(c model.CharacterClass) string {
	if label, exists := characterClassLabels[c]; exists {
		return label
	} else {
		return ""
	}
}

var reverseRaidLabels map[string]model.RaidName
var raidNameLabels = map[model.RaidName]string{
	model.SALVATIONS_EDGE:     "Salvation's Edge",
	model.CROTAS_END:          "Crota's End",
	model.ROOT_OF_NIGHTMARES:  "Root of Nightmares",
	model.KINGS_FALL:          "King's Fall",
	model.VOW_OF_THE_DISCIPLE: "Vow of the Disciple",
	model.VAULT_OF_GLASS:      "Vault of Glass",
	model.DEEP_STONE_CRYPT:    "Deep Stone Crypt",
	model.GARDEN_OF_SALVATION: "Garden of Salvation",
	model.CROWN_OF_SORROW:     "Crown of Sorrow",
	model.LAST_WISH:           "Last Wish",
	model.SPIRE_OF_STARS:      "Leviathan, Spire of Stars",
	model.EATER_OF_WORLDS:     "Leviathan, Eater of Worlds",
	model.LEVIATHAN:           "Leviathan",
	model.SCOURGE_OF_THE_PAST: "Scourge of the Past",
}

func RaidLabel(rn model.RaidName) string {
	if label, exists := raidNameLabels[rn]; exists {
		return label
	} else {
		return ""
	}
}

var reverseDifficultyLabels map[string]model.RaidDifficulty
var raidDifficultyLabels = map[model.RaidDifficulty]string{
	model.NORMAL:         "Normal",
	model.PRESTIGE:       "Prestige",
	model.MASTER:         "Master",
	model.GUIDED_GAMES:   "Guided Games",
	model.CHALLENGE_MODE: "Challenge Mode",
}

var once sync.Once

// Initializes the reverse look up maps only once
func initReverseMaps() {
	once.Do(func() {
		reverseDifficultyLabels = make(map[string]model.RaidDifficulty)
		for raidDifficulty, label := range raidDifficultyLabels {
			reverseDifficultyLabels[label] = raidDifficulty
		}

		reverseRaidLabels = make(map[string]model.RaidName)
		for raidName, label := range raidNameLabels {
			reverseRaidLabels[label] = raidName
		}

		reverseClassLabels = make(map[string]model.CharacterClass)
		for characterClass, label := range characterClassLabels {
			reverseClassLabels[label] = characterClass
		}
	})
}

// This method returns an instance of a Raid Name and RaidDifficulty
// given a string, e.g., Last Wish should yield both the RaidName LAST_WISH
// and the RaidDiffculty NORMAL. On the other hand, Salvation's Edge: Master
// should yield RaidName SALVATIONS_EDGE and RaidDifficulty MASTER
func GetRaidAndDifficulty(label string) (model.RaidName, model.RaidDifficulty, error) {
	initReverseMaps()
	tokens := strings.Split(label, ":")

	if len(tokens) <= 0 {
		log.Panicf("Unable to tokenize raid Manifest Display Name [%s]", label)
		return "", "", errors.New("Unable to tokenize raid Manifest Display Name")
	}
	name := strings.TrimSpace(tokens[0])
	raidName, nameExists := reverseRaidLabels[name]

	if !nameExists {
		return "", "", fmt.Errorf("Raid name [%s] has no match", name)
	}

	if len(tokens) <= 1 {
		return raidName, model.NORMAL, nil
	}

	difficulty := strings.TrimSpace(tokens[1]) // Default difficulty
	raidDifficulty, difficultyExists := reverseDifficultyLabels[difficulty]
	log.Printf("Difficulty exists for raid [%s]: %v", label, difficultyExists)
	if !difficultyExists {
		switch {
		case strings.EqualFold(difficulty, "Standard"):
			raidDifficulty = reverseDifficultyLabels["Normal"]
		case strings.EqualFold(difficulty, "Expert") || strings.EqualFold(difficulty, "Legend"):
			raidDifficulty = reverseDifficultyLabels["Challenge Mode"]
		default:
			return "", "", fmt.Errorf("Raid difficulty [%s] has no match", difficulty)
		}
	}

	return raidName, raidDifficulty, nil
}
