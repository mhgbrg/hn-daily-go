package models

import (
	"fmt"
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

func (date Date) ToYearMonth() YearMonth {
	return YearMonth{
		Year:  date.Year,
		Month: date.Month,
	}
}

func (date Date) Prev() Date {
	t := date.ToTime()
	t = t.AddDate(0, 0, 1)
	return FromTime(t)
}

func (date Date) Next() Date {
	t := date.ToTime()
	t = t.AddDate(0, 0, 1)
	return FromTime(t)
}

func (date Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year, date.Month, date.Day)
}

func (date *Date) Scan(src interface{}) error {
	srcTime, ok := src.(time.Time)
	if !ok {
		return errors.Errorf("failed to cast value %s to time.Time", src)
	}

	*date = FromTime(srcTime)
	return nil
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
		return Date{}, errors.Errorf("invalid date: %s", str)
	}
	y, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	d, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return Date{}, errors.Errorf("invalid date: %s", str)
	}
	return Date{Year: y, Month: time.Month(m), Day: d}, nil
}
