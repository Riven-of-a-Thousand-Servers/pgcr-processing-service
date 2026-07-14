package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"pgcr-processing-service/internal/types"
)

var (
	reverseClassLabels   map[string]types.CharacterClass
	characterClassLabels = map[types.CharacterClass]string{
		types.TITAN:   "Titan",
		types.WARLOCK: "Warlock",
		types.HUNTER:  "Hunter",
	}
)

func ClassLabel(c types.CharacterClass) string {
	if label, exists := characterClassLabels[c]; exists {
		return label
	} else {
		return ""
	}
}

var (
	reverseRaidLabels map[string]types.RaidName
	raidNameLabels    = map[types.RaidName]string{
		types.SALVATIONS_EDGE:     "Salvation's Edge",
		types.CROTAS_END:          "Crota's End",
		types.ROOT_OF_NIGHTMARES:  "Root of Nightmares",
		types.KINGS_FALL:          "King's Fall",
		types.VOW_OF_THE_DISCIPLE: "Vow of the Disciple",
		types.VAULT_OF_GLASS:      "Vault of Glass",
		types.DEEP_STONE_CRYPT:    "Deep Stone Crypt",
		types.GARDEN_OF_SALVATION: "Garden of Salvation",
		types.CROWN_OF_SORROW:     "Crown of Sorrow",
		types.LAST_WISH:           "Last Wish",
		types.SPIRE_OF_STARS:      "Leviathan, Spire of Stars",
		types.EATER_OF_WORLDS:     "Leviathan, Eater of Worlds",
		types.LEVIATHAN:           "Leviathan",
		types.SCOURGE_OF_THE_PAST: "Scourge of the Past",
	}
)

func RaidLabel(rn types.RaidName) string {
	if label, exists := raidNameLabels[rn]; exists {
		return label
	} else {
		return ""
	}
}

var (
	reverseDifficultyLabels map[string]types.RaidDifficulty
	raidDifficultyLabels    = map[types.RaidDifficulty]string{
		types.NORMAL:         "Normal",
		types.PRESTIGE:       "Prestige",
		types.MASTER:         "Master",
		types.GUIDED_GAMES:   "Guided Games",
		types.CHALLENGE_MODE: "Challenge Mode",
	}
)

var once sync.Once

// Initializes the reverse look up maps only once
func initReverseMaps() {
	once.Do(func() {
		reverseDifficultyLabels = make(map[string]types.RaidDifficulty)
		for raidDifficulty, label := range raidDifficultyLabels {
			reverseDifficultyLabels[label] = raidDifficulty
		}

		reverseRaidLabels = make(map[string]types.RaidName)
		for raidName, label := range raidNameLabels {
			reverseRaidLabels[label] = raidName
		}

		reverseClassLabels = make(map[string]types.CharacterClass)
		for characterClass, label := range characterClassLabels {
			reverseClassLabels[label] = characterClass
		}
	})
}

// This method returns an instance of a Raid Name and RaidDifficulty
// given a string, e.g., Last Wish should yield both the RaidName LAST_WISH
// and the RaidDiffculty NORMAL. On the other hand, Salvation's Edge: Master
// should yield RaidName SALVATIONS_EDGE and RaidDifficulty MASTER
func GetRaidAndDifficulty(label string) (types.RaidName, types.RaidDifficulty, error) {
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
		return raidName, types.NORMAL, nil
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

func GetDamageType(enumValue int) types.DamageType {
	switch enumValue {
	case 1:
		return types.KINETIC
	case 2:
		return types.ARC
	case 3:
		return types.SOLAR
	case 4:
		return types.VOID
	case 6:
		return types.STASIS
	case 7:
		return types.STRAND
	default:
		return ""
	}
}

type EquippingBlocktypes interface {
	~int64 | ~string | ~int
}

func GetEquippingSlot[T EquippingBlocktypes](enumValue T) types.EquipmentSlot {
	value := any(enumValue)
	if v, ok := value.(int64); ok {
		switch v {
		case 1498876634:
			return types.PRIMARY
		case 2465295065:
			return types.SPECIAL
		case 953998645:
			return types.HEAVY
		}
	} else if s, ok := value.(string); ok {
		switch strings.ToLower(s) {
		case "kinetic weapons":
			return types.PRIMARY
		case "energy weapons":
			return types.SPECIAL
		case "power weapons":
			return types.HEAVY
		}
	}
	return ""
}
