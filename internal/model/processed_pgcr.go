package model

import (
	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"time"
)

type ProcessedPostGameCarnageReport struct {
	StartTime         time.Time            `json:"startTime"`
	EndTime           time.Time            `json:"endTime"`
	FromBeginning     bool                 `json:"fromBeginning"`
	InstanceId        int64                `json:"instanceId"`
	RaidName          types.RaidName       `json:"raidName"`
	RaidDifficulty    types.RaidDifficulty `json:"raidDifficulty"`
	ActivityHash      int64                `json:"activityHash"`
	Flawless          bool                 `json:"flawless"`
	Solo              bool                 `json:"solo"`
	Duo               bool                 `json:"duo"`
	Trio              bool                 `json:"trio"`
	PlayerInformation []PlayerInformation  `json:"playerInformation"`
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
	CharacterClass     types.CharacterClass         `json:"characterClass"`
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
