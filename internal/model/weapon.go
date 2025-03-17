package model

type WeaponEntity struct {
	WeaponHash          int64
	WeaponIcon          string
	WeaponName          string
	WeaponDamageType    DamageType
	WeaponEquipmentSlot EquipmentSlot
}

type DamageType string

const (
	KINETIC DamageType = "KINETIC"
	ARC     DamageType = "ARC"
	VOID    DamageType = "VOID"
	SOLAR   DamageType = "SOLAR"
	STASIS  DamageType = "STASIS"
	STRAND  DamageType = "STRAND"
)

type EquipmentSlot string

const (
	PRIMARY EquipmentSlot = "PRIMARY"
	SPECIAL EquipmentSlot = "SPECIAL"
	HEAVY   EquipmentSlot = "HEAVY"
)
