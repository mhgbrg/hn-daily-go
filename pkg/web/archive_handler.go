package web

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
)

func Archive(templates *Templates, db *sql.DB, digestRepo repo.DigestRepo) CustomHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		yearMonthStr := mux.Vars(r)["yearMonth"]
		yearMonth, err := models.ParseYearMonth(yearMonthStr)
		if err != nil {
			return NotFoundError(err)
		}

		dates, err := digestRepo.LoadDatesWithDigests(db, yearMonth)
		if err != nil {
			return InternalServerError(err)
		}
		if len(dates) == 0 {
			return NotFoundError(errors.New("no digest for month"))
		}

		firstDigest, err := digestRepo.LoadFirstDigest(db)
		if err != nil {
			return InternalServerError(err)
		}
		lastDigest, err := digestRepo.LoadLatestDigest(db)
		if err != nil {
			return InternalServerError(err)
		}

		firstYearMonth := firstDigest.Date.ToYearMonth()
		lastYearMonth := lastDigest.Date.ToYearMonth()

		viewData := createArchiveViewData(yearMonth, dates, firstYearMonth, lastYearMonth)
		err = templates.Archive.Execute(w, viewData)
		if err != nil {
			return InternalServerError(err)
		}

		return nil
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

func createArchiveViewData(
	yearMonth models.YearMonth,
	dates []models.Date,
	firstYearMonth,
	lastYearMonth models.YearMonth,
) archiveViewData {
	viewDates := make([]archiveViewDate, len(dates))
	for i, date := range dates {
		viewDates[i] = archiveViewDate{
			Date:      date.String(),
			DigestURL: DigestURL(date),
		}
	}

	prevMonthURL := ArchiveURL(yearMonth.PrevMonth())
	nextMonthURL := ArchiveURL(yearMonth.NextMonth())
	if yearMonth == firstYearMonth {
		prevMonthURL = ""
	}
	if yearMonth == lastYearMonth {
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
