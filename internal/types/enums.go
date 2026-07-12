package types

type CharacterClass string

const (
	TITAN   CharacterClass = "Titan"
	WARLOCK CharacterClass = "Warlock"
	HUNTER  CharacterClass = "Hunter"
)

type RaidName string

const (
	SALVATIONS_EDGE     RaidName = "Salvation's Edge"
	CROTAS_END          RaidName = "Crota's End"
	ROOT_OF_NIGHTMARES  RaidName = "Root of Nightmares"
	KINGS_FALL          RaidName = "King's Fall"
	VOW_OF_THE_DISCIPLE RaidName = "Vow of the Disciple"
	VAULT_OF_GLASS      RaidName = "Vault of Glass"
	DEEP_STONE_CRYPT    RaidName = "Deep Stone Crypt"
	GARDEN_OF_SALVATION RaidName = "Garden of Salvation"
	CROWN_OF_SORROW     RaidName = "Crown of Sorrow"
	LAST_WISH           RaidName = "Last Wish"
	SPIRE_OF_STARS      RaidName = "Leviathan, Spire of Stars"
	EATER_OF_WORLDS     RaidName = "Leviathan, Eater of Worlds"
	LEVIATHAN           RaidName = "Leviathan"
	SCOURGE_OF_THE_PAST RaidName = "Scourge of the Past"
)

type RaidDifficulty string

const (
	NORMAL         RaidDifficulty = "Normal"
	PRESTIGE       RaidDifficulty = "Prestige"
	MASTER         RaidDifficulty = "Master"
	GUIDED_GAMES   RaidDifficulty = "Guided Games"
	CHALLENGE_MODE RaidDifficulty = "Challenge Mode"
)

type DamageType string

const (
	KINETIC DamageType = "Kinetic"
	ARC     DamageType = "Arc"
	VOID    DamageType = "Void"
	SOLAR   DamageType = "Solar"
	STASIS  DamageType = "Stasis"
	STRAND  DamageType = "Strand"
)

type EquipmentSlot string

const (
	PRIMARY EquipmentSlot = "Primary"
	SPECIAL EquipmentSlot = "Special"
	HEAVY   EquipmentSlot = "Heavy"
)

type ClearType string

const (
	SOLO          ClearType = "Solo"
	DUO           ClearType = "Duo"
	TRIO          ClearType = "Trio"
	SOLO_FLAWLESS ClearType = "Solo Flawless"
	DUO_FLAWLESS  ClearType = "Duo Flawless"
	TRIO_FLAWLESS ClearType = "Trio Flawless"
)
