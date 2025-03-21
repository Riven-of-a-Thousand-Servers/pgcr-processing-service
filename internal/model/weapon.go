package model

import (
	types "github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

type WeaponEntity struct {
	WeaponHash          int64
	WeaponIcon          string
	WeaponName          string
	WeaponDamageType    types.DamageType
	WeaponEquipmentSlot types.EquipmentSlot
}
