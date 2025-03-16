package model

import (
	"time"
)

type CharacterClass string

const (
	TITAN   CharacterClass = "TITAN"
	WARLOCK CharacterClass = "WARLOCK"
	HUNTER  CharacterClass = "HUNTER"
)

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

type RaidDifficulty string

const (
	NORMAL         RaidDifficulty = "NORMAL"
	PRESTIGE       RaidDifficulty = "PRESTIGE"
	MASTER         RaidDifficulty = "MASTER"
	GUIDED_GAMES   RaidDifficulty = "GUIDED_GAMES"
	CHALLENGE_MODE RaidDifficulty = "CHALLENGE_MODE"
)

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
