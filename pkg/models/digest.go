package models

import (
	"time"
)

type Digest struct {
	ID          int
	Date        Date
	GeneratedAt time.Time
	Stories     []Story
}
