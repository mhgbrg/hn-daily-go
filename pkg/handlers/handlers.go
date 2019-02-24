package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mhgbrg/hndaily/pkg"
)

func GetDigest(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	t, err := template.ParseFiles("templates/digest.html")
	if err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Path[len("/digest/"):]
		date, err := pkg.ParseDate(dateStr)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}

		digest, err := pkg.LoadDigest(db, date)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}

		viewDigest := createViewData(digest)
		var responseBody bytes.Buffer
		err = t.Execute(&responseBody, viewDigest)
		if err != nil {
			log.Printf("%+v\n", err)
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, responseBody.String())
	}
}

type ViewData struct {
	Digest ViewDigest
}

type ViewDigest struct {
	Year        int
	Month       string
	Day         int
	Weekday     string
	GeneratedAt time.Time
	Stories     []ViewStory
}

type ViewStory struct {
	Rank        int
	Title       string
	URL         string
	Site        string
	Points      int
	NumComments int
	CommentsURL string
	IsRead      bool
}

func createViewData(digest pkg.Digest) ViewData {
	viewDigest := ViewDigest{
		Weekday:     digest.Date.ToTime().Weekday().String(),
		Month:       digest.Date.Month.String(),
		Day:         digest.Date.Day,
		Year:        digest.Date.Year,
		GeneratedAt: digest.GeneratedAt,
		Stories:     make([]ViewStory, len(digest.Stories)),
	}

	for i, story := range digest.Stories {
		viewStory := ViewStory{
			Rank:        i + 1,
			Title:       story.Title,
			URL:         story.URL.String(),
			Site:        story.URL.Hostname(),
			Points:      story.Points,
			NumComments: story.NumComments,
			CommentsURL: fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ExternalID),
			IsRead:      false,
		}
		viewDigest.Stories[i] = viewStory
	}

	return ViewData{Digest: viewDigest}
}
