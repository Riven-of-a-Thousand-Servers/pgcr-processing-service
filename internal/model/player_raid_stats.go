package model

type PlayerRaidStatsEntity struct {
	RaidName           RaidName
	RaidDifficulty     RaidDifficulty
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
