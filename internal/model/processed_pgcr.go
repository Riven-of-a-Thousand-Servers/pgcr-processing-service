package model

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type CharacterClass string

const (
	TITAN   CharacterClass = "TITAN"
	WARLOCK CharacterClass = "WARLOCK"
	HUNTER  CharacterClass = "HUNTER"
)

var reverseCharacterClassLabels map[string]CharacterClass
var characterClassLabels = map[CharacterClass]string{
	TITAN:   "Titan",
	WARLOCK: "Warlock",
	HUNTER:  "Hunter",
}

func (c CharacterClass) Label() string {
	if label, exists := characterClassLabels[c]; exists {
		return label
	} else {
		return ""
	}
}

type RaidName string

const (
	SALVATIONS_EDGE     RaidName = "SALVATIONS_EDGE"
	CROTAS_END          RaidName = "CROTAS_END"
	ROOT_OF_NIGHTMARES  RaidName = "ROOT_OF_NIGHTMARES"
	KINGS_FALL          RaidName = "KINGS_FALL"
	VOW_OF_THE_DISCIPLE RaidName = "VOW_OF_THE_DISCIPLE"
	VAULT_OF_GLASS      RaidName = "VAULT_OF_GLASS"
	DEEP_STONE_CRYPT    RaidName = "DEEP_STONE_CRYPT"
	GARDEN_OF_SALVATION RaidName = "GARDEN_OF_SALVATION"
	CROWN_OF_SORROW     RaidName = "CROWN_OF_SORROW"
	LAST_WISH           RaidName = "LAST_WISH"
	SPIRE_OF_STARS      RaidName = "SPIRE_OF_STARS"
	EATER_OF_WORLDS     RaidName = "EATER_OF_WORLDS"
	LEVIATHAN           RaidName = "LEVIATHAN"
	SCOURGE_OF_THE_PAST RaidName = "SCOURGE_OF_THE_PAST"
)

var reverseRaidNameLabels map[string]RaidName
var raidNameLabels = map[RaidName]string{
	SALVATIONS_EDGE:     "Salvation's Edge",
	CROTAS_END:          "Crota's End",
	ROOT_OF_NIGHTMARES:  "Root of Nightmares",
	KINGS_FALL:          "King's Fall",
	VOW_OF_THE_DISCIPLE: "Vow of the Discple",
	VAULT_OF_GLASS:      "Vault of Glass",
	DEEP_STONE_CRYPT:    "Deep Stone Crypt",
	GARDEN_OF_SALVATION: "Garden of Salvation",
	CROWN_OF_SORROW:     "Crown of Sorrow",
	LAST_WISH:           "Last Wish",
	SPIRE_OF_STARS:      "Leviathan: Spire of Stars",
	EATER_OF_WORLDS:     "Leviathan: Eater of Worlds",
	LEVIATHAN:           "Leviathan",
	SCOURGE_OF_THE_PAST: "Scourge of the Past",
}

func (rn RaidName) Label() string {
	if label, exists := raidNameLabels[rn]; exists {
		return label
	} else {
		return ""
	}
}

type RaidDifficulty string

const (
	NORMAL       RaidDifficulty = "NORMAL"
	PRESTIGE     RaidDifficulty = "PRESTIGE"
	MASTER       RaidDifficulty = "MASTER"
	GUIDED_GAMES RaidDifficulty = "GUIDED_GAMES"
)

var reverseRaidDifficultyLabels map[string]RaidDifficulty
var raidDifficultyLabels = map[RaidDifficulty]string{
	NORMAL:       "Normal",
	PRESTIGE:     "Prestige",
	MASTER:       "Master",
	GUIDED_GAMES: "Guided Games",
}

func (rd RaidDifficulty) Label() string {
	if label, exists := raidDifficultyLabels[rd]; exists {
		return label
	} else {
		return ""
	}
}

var once sync.Once

// Initializes the reverse look up maps only once
func initReverseMaps() {
	once.Do(func() {
		reverseRaidDifficultyLabels = make(map[string]RaidDifficulty)
		for raidDifficulty, label := range raidDifficultyLabels {
			reverseRaidDifficultyLabels[label] = raidDifficulty
		}

		reverseRaidNameLabels = make(map[string]RaidName)
		for raidName, label := range raidNameLabels {
			reverseRaidNameLabels[label] = raidName
		}

		reverseCharacterClassLabels = make(map[string]CharacterClass)
		for characterClass, label := range characterClassLabels {
			reverseCharacterClassLabels[label] = characterClass
		}
	})
}

// Get a raid difficulty from a string label
func RD(label string) (RaidDifficulty, error) {
	initReverseMaps()
	if raidDifficulty, exists := reverseRaidDifficultyLabels[label]; exists {
		return raidDifficulty, nil
	} else {
		return "", errors.New("Invalid raid difficulty label")
	}
}

func Raid(label string) (RaidName, RaidDifficulty, error) {
	initReverseMaps()
	tokens := strings.Split(label, ":")

	if len(tokens) <= 0 {
		log.Panicf("Unable to tokenize raid Manifest Display Name [%s]", label)
		return "", "", errors.New("Unable to tokenize raid Manifest Display Name")
	}
	rawRaidName := strings.TrimSpace(tokens[0])
	rawRaidDifficulty := "Normal" // Default difficulty
	if len(tokens) > 1 {
		rawRaidDifficulty = strings.TrimSpace(tokens[1])
	}

	raidName, a := reverseRaidNameLabels[rawRaidName]
	raidDifficulty, b := reverseRaidDifficultyLabels[rawRaidDifficulty]
	if a && b {
		return raidName, raidDifficulty, nil
	} else {
		return "", "", fmt.Errorf("RaidName [%s] exists: [%b]. Raid difficulty [%s] exists: [%b]", raidName, a, raidDifficulty, b)
	}
}

// Get a raid name from a string label
func RN(label string) (RaidName, error) {
	initReverseMaps()
	if raidName, exists := reverseRaidNameLabels[label]; exists {
		return raidName, nil
	} else {
		return "", errors.New("Invalid raid difficulty label")
	}
}

func CC(label string) (CharacterClass, error) {
	initReverseMaps()
	if characterClass, exists := reverseCharacterClassLabels[label]; exists {
		return characterClass, nil
	} else {
		return "", errors.New("Invalid character class label")
	}
}

type ProcessedPostGameCarnageReport struct {
	StartTime         time.Time           `json:"startTime"`
	EndTime           time.Time           `json:"endTime"`
	FromBeginning     bool                `json:"fromBeginning"`
	InstanceId        int64               `json:"instanceId"`
	RaidName          RaidName            `json:"raidName"`
	RaidDifficulty    RaidDifficulty      `json:"raidDifficulty"`
	ActivityHash      int64               `json:"activityHash"`
	Flawless          bool                `json:"flawless"`
	Solo              bool                `json:"solo"`
	Duo               bool                `json:"duo"`
	Trio              bool                `json:"trio"`
	PlayerInformation []PlayerInformation `json:"playerInformation"`
}

type PlayerInformation struct {
	MembershipId               int64                        `json:"membershipId"`
	MembershipType             int                          `json:"membershipType"`
	DisplayName                string                       `json:"displayName"`
	GlobalDisplayName          string                       `json:"globalDisplayName"`
	GlobalDisplayNameCode      int                          `json:"globalDisplayNameCode"`
	PlayerCharacterInformation []PlayerCharacterInformation `json:"characterInformation"`
}

type PlayerCharacterInformation struct {
	CharacterId        int64                        `json:"characterId"`
	LightLevel         int                          `json:"lightLevel"`
	CharacterClass     CharacterClass               `json:"characterClass"`
	CharacterEmblem    int64                        `json:"characterEmblem"`
	ActivityCompleted  bool                         `json:"activityCompleted"`
	Kills              int                          `json:"kills"`
	Assists            int                          `json:"assists"`
	Deaths             int                          `json:"deaths"`
	Kda                float32                      `json:"kda"`
	Kdr                float32                      `json:"kdr"`
	TimePlayedSeconds  int                          `json:"timePlayedSeconds"`
	WeaponInformation  []CharacterWeaponInformation `json:"weaponInformation"`
	AbilityInformation CharacterAbilityInformation  `json:"abilityInformation"`
}

type CharacterWeaponInformation struct {
	WeaponHash     int64   `json:"weaponHash"`
	Kills          int     `json:"kills"`
	PrecisionKills int     `json:"precisionKills"`
	PrecisionRatio float32 `json:"precisionRatio"`
}

type CharacterAbilityInformation struct {
	GrenadeKills int `json:"grenadeKills"`
	MeleeKills   int `json:"meleeKills"`
	SuperKills   int `json:"superKills"`
}
