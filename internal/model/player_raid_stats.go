package model

import (
	types "github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
)

type PlayerRaidStatsEntity struct {
	RaidName           types.RaidName
	RaidDifficulty     types.RaidDifficulty
	PlayerMembershipId int64
	Kills              int
	Deaths             int
	Assists            int
	HoursPlayed        int
	Clears             int
	FullClears         int
	Flawless           bool
	ContestClear       bool
	DayOne             bool
	Solo               bool
	Duo                bool
	Trio               bool
	SoloFlawless       bool
	DuoFlawless        bool
	TrioFlawless       bool
}
