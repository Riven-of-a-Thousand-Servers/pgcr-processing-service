package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"pgcr-processing-service/internal/types/rabbitmq"
)

var (
	reverseClassLabels   map[string]rabbitmq.CharacterClass
	characterClassLabels = map[rabbitmq.CharacterClass]string{
		rabbitmq.TITAN:   "Titan",
		rabbitmq.WARLOCK: "Warlock",
		rabbitmq.HUNTER:  "Hunter",
	}
)

func ClassLabel(c rabbitmq.CharacterClass) string {
	if label, exists := characterClassLabels[c]; exists {
		return label
	} else {
		return ""
	}
}

var (
	reverseRaidLabels map[string]rabbitmq.RaidName
	raidNameLabels    = map[rabbitmq.RaidName]string{
		rabbitmq.SALVATIONS_EDGE:     "Salvation's Edge",
		rabbitmq.CROTAS_END:          "Crota's End",
		rabbitmq.ROOT_OF_NIGHTMARES:  "Root of Nightmares",
		rabbitmq.KINGS_FALL:          "King's Fall",
		rabbitmq.VOW_OF_THE_DISCIPLE: "Vow of the Disciple",
		rabbitmq.VAULT_OF_GLASS:      "Vault of Glass",
		rabbitmq.DEEP_STONE_CRYPT:    "Deep Stone Crypt",
		rabbitmq.GARDEN_OF_SALVATION: "Garden of Salvation",
		rabbitmq.CROWN_OF_SORROW:     "Crown of Sorrow",
		rabbitmq.LAST_WISH:           "Last Wish",
		rabbitmq.SPIRE_OF_STARS:      "Leviathan, Spire of Stars",
		rabbitmq.EATER_OF_WORLDS:     "Leviathan, Eater of Worlds",
		rabbitmq.LEVIATHAN:           "Leviathan",
		rabbitmq.SCOURGE_OF_THE_PAST: "Scourge of the Past",
	}
)

func RaidLabel(rn rabbitmq.RaidName) string {
	if label, exists := raidNameLabels[rn]; exists {
		return label
	} else {
		return ""
	}
}

var (
	reverseDifficultyLabels map[string]rabbitmq.RaidDifficulty
	raidDifficultyLabels    = map[rabbitmq.RaidDifficulty]string{
		rabbitmq.NORMAL:         "Normal",
		rabbitmq.PRESTIGE:       "Prestige",
		rabbitmq.MASTER:         "Master",
		rabbitmq.GUIDED_GAMES:   "Guided Games",
		rabbitmq.CHALLENGE_MODE: "Challenge Mode",
	}
)

var once sync.Once

// Initializes the reverse look up maps only once
func initReverseMaps() {
	once.Do(func() {
		reverseDifficultyLabels = make(map[string]rabbitmq.RaidDifficulty)
		for raidDifficulty, label := range raidDifficultyLabels {
			reverseDifficultyLabels[label] = raidDifficulty
		}

		reverseRaidLabels = make(map[string]rabbitmq.RaidName)
		for raidName, label := range raidNameLabels {
			reverseRaidLabels[label] = raidName
		}

		reverseClassLabels = make(map[string]rabbitmq.CharacterClass)
		for characterClass, label := range characterClassLabels {
			reverseClassLabels[label] = characterClass
		}
	})
}

// This method returns an instance of a Raid Name and RaidDifficulty
// given a string, e.g., Last Wish should yield both the RaidName LAST_WISH
// and the RaidDiffculty NORMAL. On the other hand, Salvation's Edge: Master
// should yield RaidName SALVATIONS_EDGE and RaidDifficulty MASTER
func GetRaidAndDifficulty(label string) (rabbitmq.RaidName, rabbitmq.RaidDifficulty, error) {
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
		return raidName, rabbitmq.NORMAL, nil
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

func GetDamageType(enumValue int) rabbitmq.DamageType {
	switch enumValue {
	case 1:
		return rabbitmq.KINETIC
	case 2:
		return rabbitmq.ARC
	case 3:
		return rabbitmq.SOLAR
	case 4:
		return rabbitmq.VOID
	case 6:
		return rabbitmq.STASIS
	case 7:
		return rabbitmq.STRAND
	default:
		return ""
	}
}

type EquippingBlocktypes interface {
	~int64 | ~string | ~int
}

func GetEquippingSlot[T EquippingBlocktypes](enumValue T) rabbitmq.EquipmentSlot {
	value := any(enumValue)
	if v, ok := value.(int64); ok {
		switch v {
		case 1498876634:
			return rabbitmq.PRIMARY
		case 2465295065:
			return rabbitmq.SPECIAL
		case 953998645:
			return rabbitmq.HEAVY
		}
	} else if s, ok := value.(string); ok {
		switch strings.ToLower(s) {
		case "kinetic weapons":
			return rabbitmq.PRIMARY
		case "energy weapons":
			return rabbitmq.SPECIAL
		case "power weapons":
			return rabbitmq.HEAVY
		}
	}
	return ""
}
