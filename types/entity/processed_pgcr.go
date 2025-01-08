package entity 

import (
	"time"
)

type CharacterClass string

const (
	Titan   CharacterClass = "TITAN"
	Warlock CharacterClass = "WARLOCK"
	Hunter  CharacterClass = "HUNTER"
)

type CharacterGender string

const (
	Male   CharacterGender = "MALE"
	Female CharacterGender = "FEMALE"
)

type CharacterRace string

const (
	Human  CharacterRace = "HUMAN"
	Awoken CharacterRace = "AWOKEN"
	Exo    CharacterRace = "EXO"
)

type RaidName string

const (
	SalvationsEdge         RaidName = "SALVATIONS_EDGE"
	CrotasEnd              RaidName = "CROTAS_END"
	RootOfNightmares       RaidName = "ROOT_OF_NIGHTMARES"
	KingsFall              RaidName = "KINGS_FALL"
	VowOfTheDisciple       RaidName = "VOW_OF_THE_DISCIPLE"
	VaultOfGlass           RaidName = "VAULT_OF_GLASS"
	DeepStoneCrypt         RaidName = "DEEP_STONE_CRYPT"
	GardenOfSalvation      RaidName = "GARDEN_OF_SALVATION"
	LeviathanCrownOfSorrow RaidName = "LEVIATHAN_CROWN_OF_SORROW"
	LastWish               RaidName = "LAST_WISH"
	LeviathanSpireOfStars  RaidName = "LEVIATHAN_SPIRE_OF_STARS"
	LeviathanEaterOfWorlds RaidName = "LEVIATHAN_EATER_OF_WORLDS"
	Leviathan              RaidName = "LEVIATHAN"
	ScourgeOfThePast       RaidName = "SCOURGE_OF_THE_PAST"
)

type RaidDifficulty string

const (
	Normal   RaidDifficulty = "NORMAL"
	Prestige RaidDifficulty = "PRESTIGE"
	Master   RaidDifficulty = "MASTER"
)

type ProcessedPostGameCarnageReport struct {
	StartTime         time.Time           `json:"startTime"`
	EndTime           time.Time           `json:"endTime"`
	FromBeginning     bool                `json:"fromBeginning"`
	InstanceId        int64               `json:"instanceId"`
	RaidName          RaidName            `json:"raidName"`
	RaidDifficulty    RaidDifficulty      `json:"raidDifficulty"`
	ActivityHash      string              `json:"activityHash"`
	Flawless          bool                `json:"flawless"`
	Solo              bool                `json:"solo"`
	Duo               bool                `json:"duo"`
	Trio              bool                `json:"trio"`
	PlayerInformation []PlayerInformation `json:"playerInformation"`
}

type PlayerInformation struct {
	MembershipId               int64                      `json:"membershipId"`
	MembershipType             int                        `json:"membershipType"`
	DisplayName                string                     `json:"displayName"`
	GlobalDisplayName          string                     `json:"globalDisplayName"`
  GlobalDisplayNameCode      int                        `json:"globalDisplayNameCode"`
	PlayerCharacterInformation []PlayerCharacterInformation `json:"characterInformation"`
}

type PlayerCharacterInformation struct {
	CharacterId        int64                        `json:"characterId"`
	LightLevel         int                          `json:"lightLevel"`
	CharacterClass     CharacterClass               `json:"characterClass"`
	CharacterGender    CharacterGender              `json:"characterGender"`
	CharacterRace      CharacterRace                `json:"characterRace"`
	CharacterEmblem    int64                        `json:"characterEmblem"`
	ActivityCompleted  bool                         `json:"activityCompleted"`
	Kills              int                          `json:"kills"`
	Assists            int                          `json:"assists"`
	Deaths             int                          `json:"deaths"`
	Kda                float32                      `json:"kda"`
	Kdr                float32                      `json:"kdr"`
	TimePlayedSeconds         int                   `json:"timePlayedSeconds"`
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
