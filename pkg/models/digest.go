package models

import (
	"time"
)

type Digest struct {
	ID          int
	Date        Date
	StartTime   time.Time
	EndTime     time.Time
	GeneratedAt time.Time
	Stories     []Story
}
