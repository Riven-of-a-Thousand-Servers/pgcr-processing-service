package model

type InstanceWeaponStats struct {
	InstanceId          int64
	PlayerCharacterId   int64
	WeaponId            int64
	TotalKills          int
	TotalPrecisionKills int
	PrecisionRatio      float32
}
