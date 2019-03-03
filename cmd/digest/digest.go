package main

import (
	"log"
	"os"

	"github.com/mhgbrg/hndaily/pkg/digester"
	"github.com/mhgbrg/hndaily/pkg/models"
	"github.com/mhgbrg/hndaily/pkg/repo"
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
		err := digester.Digest(db, date)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
