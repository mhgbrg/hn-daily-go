package main

import (
	"log"
	"os"

	"github.com/mhgbrg/hndaily/pkg"
	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"

	"github.com/pkg/errors"
)

const storiesPerDigest = 10

func main() {
	args := os.Args[1:]
	if len(args) < 1 || len(args) > 2 {
		log.Fatal("usage: ./digest <date> | <start_date> <end_date>")
	}
	if len(args) == 1 {
		dateStr := args[0]
		date, err := models.ParseDate(dateStr)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		err = digestSingleDate(date)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	} else if len(args) == 2 {
		startDateStr := args[0]
		endDateStr := args[1]
		startDate, err := models.ParseDate(startDateStr)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		endDate, err := models.ParseDate(endDateStr)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		err = digestDateRange(startDate, endDate)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
	log.Print("done!")
}

func digestSingleDate(date models.Date) error {
	return digestDateRange(date, date)
}

func digestDateRange(startDate, endDate models.Date) error {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := repo.ConnectToDB(databaseURL)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	for date := startDate; date != endDate.Next(); date = date.Next() {
		log.Printf("digesting %s\n", date)

		d, err := pkg.FetchDigest(date, storiesPerDigest)
		if err != nil {
			return errors.WithMessagef(err, "failed to fetch digest for date %s", date)
		}

		tx, err := db.Begin()
		if err != nil {
			return errors.Wrap(err, "failed to open db transaction")
		}

		if err = repo.SaveDigest(tx, d); err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				return errors.Wrap(err2, "failed to rollback transaction")
			}
			return errors.WithMessagef(err, "failed to save digest for date %s, rolled back transaction", date)
		}

		if err := tx.Commit(); err != nil {
			return errors.Wrap(err, "failed to commit transaction")
		}
	}

	return nil
}
