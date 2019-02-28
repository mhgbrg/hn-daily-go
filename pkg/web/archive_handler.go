package web

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

func Archive(db *sql.DB) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (fmt.Stringer, error) {
		yearMonthStr := r.URL.Path[len("/archive/"):]
		yearMonth, err := models.ParseYearMonth(yearMonthStr)
		if err != nil {
			return nil, NotFoundError(err)
		}

		dates, err := repo.LoadDatesWithDigests(db, yearMonth)
		if err != nil {
			return nil, InternalServerError(err)
		}
		if len(dates) == 0 {
			return nil, NotFoundError(errors.New("no digest for month"))
		}

		firstDigest, err := repo.LoadFirstDigest(db)
		if err != nil {
			return nil, InternalServerError(err)
		}
		lastDigest, err := repo.LoadLatestDigest(db)
		if err != nil {
			return nil, InternalServerError(err)
		}

		firstYearMonth := models.YearMonth{Year: firstDigest.Date.Year, Month: firstDigest.Date.Month}
		lastYearMonth := models.YearMonth{Year: lastDigest.Date.Year, Month: lastDigest.Date.Month}

		viewData := createArchiveViewData(yearMonth, dates, firstYearMonth, lastYearMonth)
		if err != nil {
			return nil, InternalServerError(err)
		}

		responseBody, err := RenderTemplate("archive", viewData)
		if err != nil {
			return nil, InternalServerError(err)
		}

		return responseBody, nil
	}
}

type archiveViewData struct {
	Year         int
	Month        string
	PrevMonthURL string
	NextMonthURL string
	Dates        []archiveViewDate
}

type archiveViewDate struct {
	Date      string
	DigestURL string
}

func createArchiveViewData(yearMonth models.YearMonth, dates []models.Date, firstYearMonth, lastYearMonth models.YearMonth) archiveViewData {
	viewDates := make([]archiveViewDate, len(dates))
	for i, date := range dates {
		viewDates[i] = archiveViewDate{date.String(), fmt.Sprintf("/digest/%s", date)}
	}
	prevMonthURL := fmt.Sprintf("/archive/%s", yearMonth.PrevMonth())
	nextMonthURL := fmt.Sprintf("/archive/%s", yearMonth.NextMonth())
	if yearMonth == firstYearMonth {
		prevMonthURL = ""
	} else if yearMonth == lastYearMonth {
		nextMonthURL = ""
	}
	return archiveViewData{
		Year:         yearMonth.Year,
		Month:        yearMonth.Month.String(),
		PrevMonthURL: prevMonthURL,
		NextMonthURL: nextMonthURL,
		Dates:        viewDates,
	}
}
