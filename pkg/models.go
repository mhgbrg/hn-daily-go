package pkg

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
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

type URL url.URL

func (u *URL) Scan(src interface{}) error {
	srcStr, ok := src.(string)
	if !ok {
		return errors.New("failed to cast value %s to string")
	}

	u2, err := url.Parse(srcStr)
	if err != nil {
		return errors.New("failed to parse string %s as url")
	}

	*u = URL(*u2)
	return nil
}

func (u *URL) String() string {
	u2 := url.URL(*u)
	return u2.String()
}

func (u *URL) Hostname() string {
	u2 := url.URL(*u)
	return u2.Hostname()
}

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

type Digest struct {
	ID          int
	Date        Date
	StartTime   time.Time
	EndTime     time.Time
	GeneratedAt time.Time
	Stories     []Story
}
