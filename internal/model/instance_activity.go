package model

type InstanceActivityEntity struct {
	InstanceId         int64
	PlayerMembershipId int64
	PlayerCharacterId  int64
	CharacterEmblem    int64
	IsCompleted        bool
	Kills              int32
	Deaths             int32
	Assists            int32
	MeleeKills         int
	SuperKills         int
	GrenadeKills       int
	KillsDeathsAssists float32
	KillsDeathsRatio   float32
	Efficiency         float32
	DurationSeconds    int32
	TimeplayedSeconds  int32
}
