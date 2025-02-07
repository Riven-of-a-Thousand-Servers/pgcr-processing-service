package model

import (
	"time"
)

type RaidEntity struct {
	RaidName       string
	RaidDifficulty string
	IsActive       bool
	ReleaseDate    time.Time
}
