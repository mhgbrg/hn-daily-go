package models

import "time"

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

type Story struct {
	ID          int
	PostedAt    time.Time
	Title       string
	URL         string
	Author      string
	Points      int
	NumComments int
}

type Digest struct {
	Date        Date
	StartTime   time.Time
	EndTime     time.Time
	GeneratedAt time.Time
	Stories     []Story
}
