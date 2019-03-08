package models

import "time"

type User struct {
	ID         int
	ExternalID string
	FirstVisit time.Time
}
