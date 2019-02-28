package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type YearMonth struct {
	Year  int
	Month time.Month
}

func (yearMonth YearMonth) PrevMonth() YearMonth {
	year := yearMonth.Year
	month := yearMonth.Month - 1
	if month == 0 {
		month = 12
		year--
	}
	return YearMonth{year, month}
}

func (yearMonth YearMonth) NextMonth() YearMonth {
	year := yearMonth.Year
	month := yearMonth.Month + 1
	if month == 13 {
		month = 1
		year++
	}
	return YearMonth{year, month}
}

func (yearMonth YearMonth) String() string {
	return fmt.Sprintf("%d-%02d", yearMonth.Year, int(yearMonth.Month))
}

func ParseYearMonth(str string) (YearMonth, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 2 {
		return YearMonth{}, errors.Errorf("failed to parse string \"%s\" as YearMonth, splitting on \"-\" gave %d strings", str, len(parts))
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return YearMonth{}, errors.Errorf("failed to parse string \"%s\" as YearMonth, invalid year \"%s\"", str, parts[0])
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return YearMonth{}, errors.Errorf("failed to parse string \"%s\" as YearMonth, invalid month \"%s\"", str, parts[1])
	}

	return YearMonth{year, time.Month(month)}, nil
}
