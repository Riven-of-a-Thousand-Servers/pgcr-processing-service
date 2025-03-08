package model

import "time"

type PlayerEntity struct {
	MembershipId    int64
	DisplayName     string
	DisplayNameCode int32
	MembershipType  int32
	LastSeen        time.Time
	Characters      []PlayerCharacterEntity
}

type PlayerCharacterEntity struct {
	CharacterId        int64
	CharacterClass     string
	CharacterEmblem    string
	PlayerMembershipId int64
}
