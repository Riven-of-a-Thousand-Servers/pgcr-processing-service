package model

import (
	"time"
)

type RaidPgcr struct {
	InstanceId int64
	Timestamp  time.Time
	Blob       []byte
}
