package models

import "time"

type Story struct {
	ID          int
	ExternalID  int
	PostedAt    time.Time
	Title       string
	URL         URL
	Author      string
	Points      int
	NumComments int
}
