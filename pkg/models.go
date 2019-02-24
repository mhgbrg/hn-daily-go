package pkg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (date Date) ToTime() time.Time {
	return time.Date(date.Year, date.Month, date.Day, 0, 0, 0, 0, time.UTC)
}

func (date Date) Next() Date {
	t := date.ToTime()
	t = t.AddDate(0, 0, 1)
	return FromTime(t)
}

func (date Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year, date.Month, date.Day)
}

func FromTime(t time.Time) Date {
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

func ParseDate(str string) (Date, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 3 {
		return Date{}, fmt.Errorf("invalid date: %s", str)
	}
	y, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	d, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return Date{}, fmt.Errorf("invalid date: %s", str)
	}
	return Date{Year: y, Month: time.Month(m), Day: d}, nil
}

type Story struct {
	ExternalID  int
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
